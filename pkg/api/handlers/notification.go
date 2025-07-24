package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/r1i2t3/agni/pkg/db"
	"github.com/r1i2t3/agni/pkg/notification"
	"github.com/r1i2t3/agni/pkg/queue"
)

type NotificationRequest struct {
	Channel            notification.NotificationChannel `json:"channel"`
	Provider           string                           `json:"provider,omitempty"`
	Recipient          string                           `json:"recipient"`
	Subject            string                           `json:"subject,omitempty"`
	Message            string                           `json:"message"`
	MessageContentType string                           `json:"message_content_type,omitempty"`
	TemplateID         string                           `json:"template_id,omitempty"`
}

func EnqueueNotification(c *fiber.Ctx) error {
	// Get the application from context (stored by APIKeyAuth middleware)
	app, ok := c.Locals("app").(*db.Application)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid application context",
		})
	}

	log.Printf("API Request from App: %s (ID: %s, API Key: %s)",
		app.Name, app.ID.String(), app.APIToken)

	var request NotificationRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Use the application data
	notification := notification.Notification{
		ID:                 notification.GenerateID(),
		ApplicationID:      app.ID.String(),
		Provider:           request.Provider,
		Channel:            request.Channel,
		Recipient:          request.Recipient,
		Message:            request.Message,
		MessageContentType: request.MessageContentType,
		Subject:            request.Subject,
		Status:             "queued",
		CreatedAt:          time.Now(),
	}
	QueueID, err := queue.EnqueueNotification(notification)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to enqueue notification"})
	}

	return c.JSON(fiber.Map{"message": "Notification queued successfully", "queue_id": QueueID, "status": "queued"})
}
