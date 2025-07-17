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

// RequeueNotification puts a notification back into the queue
// Useful for retry logic or failed processing
func RequeueNotification(queueName string, notification *QueuedNotification) error {
	RedisClient := db.GetRedisClient()
	if RedisClient == nil {
		return fmt.Errorf("redis client not available")
	}

	ctx := context.Background()

	// Update the queued timestamp
	notification.QueuedAt = time.Now()

	// Serialize the notification
	data, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to serialize notification for requeue: %w", err)
	}

	// Put back at the front of the queue (LPUSH)
	err = RedisClient.LPush(ctx, queueName, data).Err()
	if err != nil {
		return fmt.Errorf("failed to requeue notification: %w", err)
	}

	log.Printf("🔄 Requeued notification %s (attempt %d)", notification.ID, notification.Attempts)
	return nil
}

// GetQueueLength returns the number of notifications in the queue
func GetQueueLength(queueName string) (int64, error) {
	RedisClient := db.GetRedisClient()
	if RedisClient == nil {
		return 0, fmt.Errorf("redis client not available")
	}

	ctx := context.Background()
	length, err := RedisClient.LLen(ctx, queueName).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get queue length: %w", err)
	}

	return length, nil
}

// PeekQueue returns the next notification without removing it from the queue
func PeekQueue(queueName string) (*QueuedNotification, error) {
	RedisClient := db.GetRedisClient()
	if RedisClient == nil {
		return nil, fmt.Errorf("redis client not available")
	}

	ctx := context.Background()

	// Use LINDEX to peek at the last item (next to be dequeued)
	result, err := RedisClient.LIndex(ctx, queueName, -1).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil // Queue is empty
		}
		return nil, fmt.Errorf("failed to peek queue: %w", err)
	}

	var queuedNotification QueuedNotification
	if err := json.Unmarshal([]byte(result), &queuedNotification); err != nil {
		return nil, fmt.Errorf("failed to parse peeked notification: %w", err)
	}

	return &queuedNotification, nil
}

// ClearQueue removes all notifications from the queue
func ClearQueue(queueName string) error {
	RedisClient := db.GetRedisClient()
	if RedisClient == nil {
		return fmt.Errorf("redis client not available")
	}

	ctx := context.Background()

	err := RedisClient.Del(ctx, queueName).Err()
	if err != nil {
		return fmt.Errorf("failed to clear queue: %w", err)
	}

	log.Printf("🗑️ Cleared queue %s", queueName)
	return nil
}
