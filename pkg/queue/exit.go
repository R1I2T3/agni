package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/r1i2t3/agni/pkg/db"
)

// DequeueResult represents the result of a dequeue operation
type DequeueResult struct {
	Notification *QueuedNotification
	Error        error
	IsEmpty      bool
}

// DequeueNotification dequeues a single notification from Redis
// Uses blocking pop with timeout to efficiently wait for notifications
func DequeueNotification(queueName string, timeout time.Duration) (*QueuedNotification, error) {
	RedisClient := db.GetRedisClient()
	if RedisClient == nil {
		return nil, fmt.Errorf("redis client not available")
	}

	ctx := context.Background()

	// Use BRPOP for blocking right pop with timeout
	result, err := RedisClient.BRPop(ctx, timeout, queueName).Result()
	if err != nil {
		// Check if it's a timeout (no items in queue)
		if err.Error() == "redis: nil" {
			return nil, fmt.Errorf("queue empty: no notifications available")
		}
		return nil, fmt.Errorf("failed to dequeue notification: %w", err)
	}

	// BRPOP returns [queueName, value]
	if len(result) < 2 {
		return nil, fmt.Errorf("invalid dequeue result format")
	}

	// Parse the notification
	var queuedNotification QueuedNotification
	if err := json.Unmarshal([]byte(result[1]), &queuedNotification); err != nil {
		return nil, fmt.Errorf("failed to parse notification: %w", err)
	}

	log.Printf("📥 Dequeued notification %s for %s", queuedNotification.ID, queuedNotification.Recipient)
	return &queuedNotification, nil
}

// DequeueNotificationNonBlocking dequeues a notification without blocking
// Returns nil if queue is empty
func DequeueNotificationNonBlocking(queueName string) (*QueuedNotification, error) {
	RedisClient := db.GetRedisClient()
	if RedisClient == nil {
		return nil, fmt.Errorf("redis client not available")
	}

	ctx := context.Background()

	// Use RPOP for non-blocking right pop
	result, err := RedisClient.RPop(ctx, queueName).Result()
	if err != nil {
		// Check if queue is empty
		if err.Error() == "redis: nil" {
			return nil, nil // Queue is empty, not an error
		}
		return nil, fmt.Errorf("failed to dequeue notification: %w", err)
	}

	// Parse the notification
	var queuedNotification QueuedNotification
	if err := json.Unmarshal([]byte(result), &queuedNotification); err != nil {
		return nil, fmt.Errorf("failed to parse notification: %w", err)
	}

	log.Printf("📥 Dequeued notification %s for %s", queuedNotification.ID, queuedNotification.Recipient)
	return &queuedNotification, nil
}
