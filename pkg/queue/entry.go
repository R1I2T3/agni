package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/r1i2t3/agni/pkg/db"
	"github.com/r1i2t3/agni/pkg/notification"
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
	err = RedisClient.LPush(ctx, "NotificationQueue", data).Err()
	if err != nil {
		return "", fmt.Errorf("failed to enqueue notification: %w", err)
	}
	log.Printf("✅ Notification queued successfully: %s", Notification.ID)
	return QueueID, nil
}
