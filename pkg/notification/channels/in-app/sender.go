package inapp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/r1i2t3/agni/pkg/db"
	"github.com/r1i2t3/agni/pkg/notification"
	"github.com/r1i2t3/agni/pkg/queue"
	"github.com/redis/go-redis/v9"
)

// InAppNotifier publishes messages to Redis Stream
type InAppNotifier struct {
	rdb    *redis.Client
	stream string
}

var InAppChannel *InAppNotifier

func NewInAppNotifier(streamName string) {
	InAppChannel = &InAppNotifier{rdb: db.GetRedisClient(), stream: streamName}
}
func (n *InAppNotifier) Send(notify *notification.Notification) error {
	b, err := json.Marshal(notify)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	args := &redis.XAddArgs{
		Stream: n.stream,

		Values: map[string]interface{}{
			"payload": string(b),
			"id":      notify.ID,
		},
	}

	if err := n.rdb.XAdd(context.Background(), args).Err(); err != nil {
		return fmt.Errorf("xadd stream %s: %w", n.stream, err)
	}
	fmt.Println("InApp notification sent")
	return nil
}

func ProcessInAppNotifications(notif *queue.QueuedNotification) (*notification.Notification, error) {
	log.Printf("Processing Inapp notification %s for %s: %+v", notif.ID, notif.Recipient, notif)
	notification := &notification.Notification{
		ID:                 notif.ID,
		ApplicationID:      notif.ApplicationID,
		QueueID:            notif.QueueID,
		Recipient:          notif.Recipient,
		Subject:            notif.Subject,
		Message:            notif.Message,
		Channel:            notif.Channel,
		Provider:           notif.Provider,
		Status:             notif.Status,
		CreatedAt:          notif.CreatedAt,
		MessageContentType: notif.MessageContentType,
		TemplateID:         notif.TemplateID,
	}
	switch notif.Provider {
	case "InApp":
		if InAppChannel == nil {
			panic("InApp notifier is not initialized")
		}
		err := InAppChannel.Send(notification)
		if err != nil {
			log.Printf("Error sending Notification: %v", err)
		}
		return notification, err

	// more providers can be added here
	default:
		log.Printf("Unknown InApp provider: %s", notif.Provider)
		return notification, fmt.Errorf("unknown InApp provider: %s", notif.Provider)
	}
}
