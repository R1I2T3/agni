package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/r1i2t3/agni/pkg/db"
	"github.com/r1i2t3/agni/pkg/queue"
)

// GetBackoffDelay calculates the exponential backoff delay based on the attempt number.
func GetBackoffDelay(attempt int) time.Duration {
	const (
		baseDelay = 100 * time.Millisecond
		maxDelay  = 10 * time.Second
	)
	delay := baseDelay * (1 << attempt)
	if delay > maxDelay {
		return maxDelay
	}
	return delay
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
func PeekQueue(queueName string) (*queue.QueuedNotification, error) {
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

	var queuedNotification queue.QueuedNotification
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
