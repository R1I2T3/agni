package sms

import (
	"context"

	"github.com/r1i2t3/agni/pkg/notification"
	"github.com/twilio/twilio-go"
)

type SMSNotifier struct {
	client     *twilio.RestClient
	fromNumber string
}

func NewSMSNotifier(accountSid, authToken, fromNumber string) *SMSNotifier {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	return &SMSNotifier{
		client:     client,
		fromNumber: fromNumber,
	}
}

func (n *SMSNotifier) Send(ctx context.Context, notification *notification.Notification) error {
	// Implement SMS sending logic using Twilio or any other SMS service
	return nil
}
