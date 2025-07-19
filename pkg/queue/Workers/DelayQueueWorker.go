package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/r1i2t3/agni/pkg/db"
	"github.com/r1i2t3/agni/pkg/queue"
	"github.com/redis/go-redis/v9"
)

// DelayedQueueProcessor handles moving ready notifications from delayed queue to main queue
type DelayedQueueProcessor struct {
	ctx              context.Context
	cancel           context.CancelFunc
	delayedQueueName string
	mainQueueName    string
	checkInterval    time.Duration
}

// NewDelayedQueueProcessor creates a new delayed queue processor
func NewDelayedQueueProcessor(delayedQueueName, mainQueueName string, checkInterval time.Duration) *DelayedQueueProcessor {
	ctx, cancel := context.WithCancel(context.Background())

	return &DelayedQueueProcessor{
		ctx:              ctx,
		cancel:           cancel,
		delayedQueueName: delayedQueueName,
		mainQueueName:    mainQueueName,
		checkInterval:    checkInterval,
	}
}

// Start begins processing delayed notifications
func (dqp *DelayedQueueProcessor) Start() {
	log.Printf("⏰ Starting delayed queue processor (checking every %v)", dqp.checkInterval)

	go func() {
		ticker := time.NewTicker(dqp.checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-dqp.ctx.Done():
				log.Println("🛑 Delayed queue processor stopping...")
				return
			case <-ticker.C:
				if err := dqp.processReadyNotifications(); err != nil {
					log.Printf("❌ Error processing delayed notifications: %v", err)
				}
			}
		}
	}()
}

// Stop gracefully stops the delayed queue processor
func (dqp *DelayedQueueProcessor) Stop() {
	dqp.cancel()
}

// processReadyNotifications moves ready notifications from delayed queue to main queue
func (dqp *DelayedQueueProcessor) processReadyNotifications() error {
	RedisClient := db.GetRedisClient()
	if RedisClient == nil {
		return fmt.Errorf("redis client not available")
	}

	now := time.Now().Unix()

	// Create a ZRangeBy struct and pass its address
	zrangeBy := &redis.ZRangeBy{
		Min: "0",
		Max: fmt.Sprintf("%d", now),
	}
	result, err := RedisClient.ZRangeByScore(dqp.ctx, dqp.delayedQueueName, zrangeBy).Result()
	if err != nil {
		return fmt.Errorf("failed to get delayed notifications: %w", err)
	}

	if len(result) == 0 {
		return nil // No notifications ready
	}

	movedCount := 0

	// Move each ready notification to the main queue
	for _, data := range result {
		// Parse notification to update status
		var queuedNotification queue.QueuedNotification
		if err := json.Unmarshal([]byte(data), &queuedNotification); err != nil {
			log.Printf("❌ Failed to parse delayed notification: %v", err)
			continue
		}

		// Update status to queued
		queuedNotification.Status = "queued"
		queuedNotification.QueuedAt = time.Now()

		// Re-serialize with updated status
		updatedData, err := json.Marshal(queuedNotification)
		if err != nil {
			log.Printf("❌ Failed to serialize updated notification: %v", err)
			continue
		}

		// Add to main queue
		err = RedisClient.LPush(dqp.ctx, dqp.mainQueueName, updatedData).Err()
		if err != nil {
			log.Printf("❌ Failed to move notification to main queue: %v", err)
			continue
		}

		// Remove from delayed queue
		err = RedisClient.ZRem(dqp.ctx, dqp.delayedQueueName, data).Err()
		if err != nil {
			log.Printf("❌ Failed to remove from delayed queue: %v", err)
		}

		movedCount++
		log.Printf("⏰ Moved delayed notification %s to main queue", queuedNotification.ID)
	}

	if movedCount > 0 {
		log.Printf("⏰ Moved %d delayed notifications to main queue", movedCount)
	}

	return nil
}

// GetDelayedQueueLength returns the number of delayed notifications
func GetDelayedQueueLength(delayedQueueName string) (int64, error) {
	RedisClient := db.GetRedisClient()
	if RedisClient == nil {
		return 0, fmt.Errorf("redis client not available")
	}

	length, err := RedisClient.ZCard(context.Background(), delayedQueueName).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get delayed queue length: %w", err)
	}

	return length, nil
}
