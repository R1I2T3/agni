package webpush

import (
	"fmt"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/r1i2t3/agni/pkg/db"
	"github.com/r1i2t3/agni/pkg/notification"
	"github.com/r1i2t3/agni/pkg/queue"
)

type PushNotifier struct {
	vapidPublicKey  string
	vapidPrivateKey string
	vapidSubject    string
}

var pushNotifier *PushNotifier

func NewPushNotifier(publicKey, privateKey, subject string) (*PushNotifier, error) {
	if publicKey == "" || privateKey == "" {
		return nil, fmt.Errorf("VAPID public and private keys are required")
	}
	pushNotifier = &PushNotifier{
		vapidPublicKey:  publicKey,
		vapidPrivateKey: privateKey,
		vapidSubject:    subject,
	}
	return pushNotifier, nil
}

func ProcessWebPushNotifications(notif *queue.QueuedNotification) (*notification.Notification, error) {
	subscriptions, err := db.GetSubscriptionByUserId(notif.Recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}
	if len(subscriptions) == 0 {
		return nil, fmt.Errorf("no subscriptions found for user: %s", notif.Recipient)
	}
	vapid := webpush.Options{
		VAPIDPublicKey:  pushNotifier.vapidPublicKey,
		VAPIDPrivateKey: pushNotifier.vapidPrivateKey,
	}
	for _, sub := range subscriptions {
		subscription := &webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				P256dh: sub.P256dh,
				Auth:   sub.Auth,
			},
		}
		resp, err := webpush.SendNotification([]byte(notif.Message), subscription, &vapid)
		if err != nil {
			return nil, fmt.Errorf("failed to send web push notification: %w", err)
		}
		resp.Body.Close()
	}
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
	return notification, nil
}
