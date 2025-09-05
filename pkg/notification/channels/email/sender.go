package email

import (
	"context"
	"fmt"
	"log"
	"net/smtp"

	"github.com/r1i2t3/agni/pkg/notification"
	"github.com/r1i2t3/agni/pkg/notification/channels/email/EmailProviders"
	"github.com/r1i2t3/agni/pkg/queue"
)

type EmailNotifier struct {
	host         string
	port         string
	username     string
	app_password string
}

var EmailChannel *EmailNotifier

func NewEmailNotifier(host, port, username, password string) error {
	if host == "" || port == "" || username == "" || password == "" {
		return fmt.Errorf("SMTP host, port, username, and password are required")
	}
	EmailChannel = &EmailNotifier{
		host:         host,
		port:         port,
		username:     username,
		app_password: password,
	}
	return nil
}

func (n *EmailNotifier) Send(ctx context.Context, notification *notification.Notification) error {
	auth := smtp.PlainAuth("", n.username, n.app_password, n.host)
	log.Default().Printf("Sending email to %s with subject %s", notification.Recipient, notification.Subject)
	to := []string{notification.Recipient}
	msg := []byte("To: " + notification.Recipient + "\r\n" +
		"Subject: " + notification.Subject + "\r\n" +
		"Content-Type: " + notification.MessageContentType + "\r\n" +
		"\r\n" + notification.Message + "\r\n")
	err := smtp.SendMail(n.host+":"+n.port, auth, n.username, to, msg)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", notification.Recipient, err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	log.Printf("Email sent to %s with subject %s\n", notification.Recipient, notification.Subject)
	notification.Status = "sent"
	return nil
}

func GetEmailChannel() *EmailNotifier {
	if EmailChannel == nil {
		panic("Email notifier is not initialized")
	}
	return EmailChannel
}

func ProcessEmailNotifications(notif *queue.QueuedNotification) (*notification.Notification, error) {
	log.Printf("Processing email notification %s for %s: %+v", notif.ID, notif.Recipient, notif)
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
	case "smtp", "email":
		if EmailChannel == nil {
			panic("Email notifier is not initialized")
		}
		err := EmailChannel.Send(context.Background(), notification)
		if err != nil {
			log.Printf("Error sending email: %v", err)
		}
		return notification, err
	case "Resend": //
		_, err := EmailProviders.Client.Send(context.Background(), notification)
		return notification, err

	// more providers can be added here
	default:
		log.Printf("Unknown email provider: %s", notif.Provider)
		return notification, fmt.Errorf("unknown email provider: %s", notif.Provider)
	}

}
