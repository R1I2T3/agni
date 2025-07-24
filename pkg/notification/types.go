package notification

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// NotificationChannel represents the delivery channel
type NotificationChannel string

const (
	ChannelEmail   NotificationChannel = "email"
	ChannelSMS     NotificationChannel = "sms"
	ChannelPush    NotificationChannel = "push"
	ChannelWebhook NotificationChannel = "webhook"
)

// IsValidChannel checks if the channel is one of the 4 allowed types
func IsValidChannel(channel string) bool {
	validChannels := []string{"email", "sms", "push", "webhook"}
	for _, valid := range validChannels {
		if channel == valid {
			return true
		}
	}
	return false
}

// ValidateChannel returns an error if channel is not valid
func ValidateChannel(channel string) error {
	if !IsValidChannel(channel) {
		return errors.New("channel must be one of: email, sms, push, webhook")
	}
	return nil
}

type Notification struct {
	ID                 string              `json:"id"`
	ApplicationID      string              `json:"application_id"`
	QueueID            string              `gorm:"type:text;uniqueIndex" json:"queue_id"`
	Channel            NotificationChannel `json:"channel" validate:"required,oneof=email sms push webhook"`
	Provider           string              `json:"provider"`
	Recipient          string              `json:"recipient" validate:"required"`
	Subject            string              `json:"subject,omitempty"`
	TemplateID         string              `json:"template_id,omitempty"`
	Message            string              `json:"message" validate:"required"`
	MessageContentType string              `json:"message_content_type,omitempty"`
	Status             string              `json:"status"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
	Attempts           int                 `json:"attempts"`
}

// SetChannel sets the channel with validation
func (n *Notification) SetChannel(channel string) error {
	if err := ValidateChannel(channel); err != nil {
		return err
	}
	n.Channel = NotificationChannel(channel)
	return nil
}

func GenerateID() string {
	//Temporary ID generation using UUID
	return uuid.New().String()
}

type Notifier interface {
	Send(ctx context.Context, notification *Notification) error
}
