package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/r1i2t3/agni/pkg/db"
	"github.com/redis/go-redis/v9"
)

func ReEnqueueNotification(QueuedNotification *QueuedNotification) (string, error) {
	RedisClient := db.GetRedisClient()
	QueueID := fmt.Sprintf("%s:%s:%s", QueuedNotification.ApplicationID, QueuedNotification.ID, QueuedNotification.Channel)

	QueuedNotification.QueueID = QueueID
	QueuedNotification.Status = "queued"
	QueuedNotification.QueuedAt = time.Now()

	// Serialize the notification
	data, err := json.Marshal(QueuedNotification)
	if err != nil {
		return "", fmt.Errorf("failed to serialize notification: %w", err)
	}
	// enqueue the notification in Redis
	err = RedisClient.LPush(ctx, "NotificationQueue", data).Err()
	if err != nil {
		return "", fmt.Errorf("failed to enqueue notification: %w", err)
	}
	log.Printf("✅ Notification queued successfully: %s", QueuedNotification.ID)
	return QueueID, nil
}

func DelayReEnqueueNotification(QueuedNotification *QueuedNotification, delay time.Duration) (string, error) {
	RedisClient := db.GetRedisClient()
	QueueID := fmt.Sprintf("%s:%s:%s", QueuedNotification.ApplicationID, QueuedNotification.ID, QueuedNotification.Channel)

	QueuedNotification.QueueID = QueueID
	QueuedNotification.Status = "scheduled"
	QueuedNotification.QueuedAt = time.Now().Add(delay)

	// Serialize the notification
	data, err := json.Marshal(QueuedNotification)
	if err != nil {
		return "", fmt.Errorf("failed to serialize notification: %w", err)
	}

	// Generate a score based on current time + delay
	score := float64(time.Now().Add(delay).Unix())
	member := redis.Z{
		Score:  score,
		Member: data,
	}
	// Add to delayed queue with score
	err = RedisClient.ZAdd(ctx, "QueuedNotification:delayed", member).Err()
	if err != nil {
		return "", fmt.Errorf("failed to enqueue delayed notification: %w", err)
	}

	log.Printf("✅ Notification queued for delayed processing: %s", QueuedNotification.ID)
	return QueueID, nil
}
