package webpush

import (
	"fmt"

	"github.com/r1i2t3/agni/pkg/notification"
	"github.com/r1i2t3/agni/pkg/queue"
)

type PushNotifier struct {
	vapidPublicKey  string
	vapidPrivateKey string
	vapidSubject    string
}

func NewPushNotifier(publicKey, privateKey, subject string) (*PushNotifier, error) {
	if publicKey == "" || privateKey == "" {
		return nil, fmt.Errorf("VAPID public and private keys are required")
	}
	return &PushNotifier{
		vapidPublicKey:  publicKey,
		vapidPrivateKey: privateKey,
		vapidSubject:    subject,
	}, nil
}

func ProcessWebPushNotifications(notif *queue.QueuedNotification) (*notification.Notification, error) {
	return nil, nil
}
