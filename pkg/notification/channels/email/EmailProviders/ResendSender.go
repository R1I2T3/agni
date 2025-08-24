package EmailProviders

import (
	"context"
	"fmt"

	"github.com/r1i2t3/agni/pkg/notification"
	"github.com/resend/resend-go/v2"
)

type ResendSender struct {
	Client *resend.Client
	from   string
}

var Client *ResendSender

func NewResendNotifier(apiKey string, from string) error {
	if apiKey == "" {
		return fmt.Errorf("resend API key is required")
	}
	Client = &ResendSender{
		Client: resend.NewClient(apiKey),
		from:   from,
	}
	return nil
}

func (s *ResendSender) Send(ctx context.Context, notification *notification.Notification) (string, error) {
	if Client == nil {
		return "", fmt.Errorf("resend notifier is not initialized")
	}

	email := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{notification.Recipient},
		Subject: notification.Subject,
	}
	switch notification.MessageContentType {
	case "text/plain":
		email.Text = notification.Message
	case "text/html":
		email.Html = notification.Message
	default:
		// handle other content types if needed
	}

	resp, err := s.Client.Emails.Send(email)
	if err != nil {
		return "", fmt.Errorf("failed to send email via Resend: %w", err)
	}

	fmt.Printf("Email sent via Resend to %s with subject %s\n", notification.Recipient, notification.Subject)

	return resp.Id, nil
}
