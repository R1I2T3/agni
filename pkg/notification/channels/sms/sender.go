package sms

import (
	"log"

	"github.com/r1i2t3/agni/pkg/notification"
	smsproviders "github.com/r1i2t3/agni/pkg/notification/channels/sms/SMSProviders"
	"github.com/r1i2t3/agni/pkg/queue"
)

func ProcessSMSNotifications(notif *queue.QueuedNotification) (*notification.Notification, error) {
	log.Printf("Processing SMS notification %s for %s: %+v", notif.ID, notif.Recipient, notif)
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
	case "twilio":
		// Process Twilio SMS notification
		_, err := smsproviders.TwilioClient.TwilioSend(notif.Recipient, notif.Message)
		if err != nil {
			log.Printf("Failed to send SMS notification via Twilio: %v", err)
		}
	default:
		log.Printf("Unknown SMS provider %s", notif.Provider)
	}
	return notification, nil
}
