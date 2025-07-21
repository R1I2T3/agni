package email

import (
	"context"
	"fmt"

	"github.com/r1i2t3/agni/pkg/notification"
)

type EmailNotifier struct {
	host         string
	port         string
	username     string
	app_password string
}

func NewEmailNotifier(host, port, username, password string) (*EmailNotifier, error) {
	if host == "" || port == "" || username == "" || password == "" {
		return nil, fmt.Errorf("SMTP host, port, username, and password are required")
	}
	return &EmailNotifier{
		host:         host,
		port:         port,
		username:     username,
		app_password: password,
	}, nil
}

func (n *EmailNotifier) Send(ctx context.Context, notification *notification.Notification) error {
	return nil
}
