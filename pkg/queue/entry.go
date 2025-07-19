package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/r1i2t3/agni/pkg/db"
	"github.com/r1i2t3/agni/pkg/notification"
	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
)

type QueuedNotification struct {
	ID            string                           `json:"id"`
	ApplicationID string                           `json:"application_id"`
	QueueID       string                           `json:"queue_id"`
	Channel       notification.NotificationChannel `json:"channel"`
	Provider      string                           `json:"provider,omitempty"`
	Recipient     string                           `json:"recipient"`
	Subject       string                           `json:"subject,omitempty"`
	Message       string                           `json:"message"`
	TemplateID    string                           `json:"template_id,omitempty"`
	Status        string                           `json:"status"`
	Attempts      int                              `json:"attempts"`
	CreatedAt     time.Time                        `json:"created_at"`
	QueuedAt      time.Time                        `json:"queued_at"`
}

func EnqueueNotification(Notification notification.Notification) (string, error) {
	RedisClient := db.GetRedisClient()
	QueueID := fmt.Sprintf("%s:%s:%s", Notification.ApplicationID, Notification.ID, Notification.Channel)
	queuedNotification := QueuedNotification{
		ID:            Notification.ID,
		ApplicationID: Notification.ApplicationID,
		QueueID:       QueueID,
		Channel:       Notification.Channel,
		Provider:      Notification.Provider,
		Recipient:     Notification.Recipient,
		Subject:       Notification.Subject,
		Message:       Notification.Message,
		TemplateID:    Notification.TemplateID,
		Status:        Notification.Status,
		Attempts:      Notification.Attempts,
		CreatedAt:     Notification.CreatedAt,
		QueuedAt:      time.Now(),
	}
	// Serialize the notification
	data, err := json.Marshal(queuedNotification)
	if err != nil {
		return "", fmt.Errorf("failed to serialize notification: %w", err)
	}
	// enqueue the notification in Redis
	err = RedisClient.LPush(ctx, "QueuedNotification", data).Err()
	if err != nil {
		return "", fmt.Errorf("failed to enqueue notification: %w", err)
	}
	log.Printf("✅ Notification queued successfully: %s", Notification.ID)
	return QueueID, nil
}

func DelayEnqueueNotification(QueuedNotification *QueuedNotification, delay time.Duration) (string, error) {
	RedisClient := db.GetRedisClient()
	if RedisClient == nil {
		return "", fmt.Errorf("redis client not available")
	}
	// Generate a unique queue ID
	QueueID := fmt.Sprintf("%s:%s:%s", QueuedNotification.ApplicationID, QueuedNotification.ID, QueuedNotification.Channel)
	QueuedNotification.QueueID = QueueID
	QueuedNotification.Status = "scheduled"
	QueuedNotification.QueuedAt = time.Now()
	QueuedNotification.CreatedAt = time.Now()

	// Serialize the notification
	data, err := json.Marshal(QueuedNotification)
	if err != nil {
		return "", fmt.Errorf("failed to serialize notification: %w", err)
	}
	//generate a score based on current time + delay
	score := float64(time.Now().Add(delay).Unix())

	member := redis.Z{
		Score:  score,
		Member: data,
	}

	// Add to delayed queue with score as current time + delay
	err = RedisClient.ZAdd(ctx, "NotificationQueue:delayed", member).Err()

	if err != nil {
		return "", fmt.Errorf("failed to enqueue delayed notification: %w", err)
	}

	log.Printf("✅ Notification queued for delayed processing: %s", QueuedNotification.ID)
	return QueueID, nil
}
