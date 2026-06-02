package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/r1i2t3/agni/pkg/db"
)

type Subscription struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		Auth   string `json:"auth"`
		P256dh string `json:"p256dh"`
	} `json:"keys"`
	Device string `json:"device,omitempty"`
	UserID string `json:"user_id,omitempty"`
}

func HandleWebPushSubscription(c *fiber.Ctx) error {
	var sub Subscription
	if err := c.BodyParser(&sub); err != nil {
		log.Printf("Error parsing body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validate required fields
	if sub.Endpoint == "" || sub.Keys.Auth == "" || sub.Keys.P256dh == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Endpoint, Auth key, and P256dh key are required"})
	}

	var subScriptionModel = &db.WebPushSubscription{
		Endpoint: sub.Endpoint,
		Auth:     sub.Keys.Auth,
		P256dh:   sub.Keys.P256dh,
		Device:   sub.Device,
		UserID:   sub.UserID,
	}
	dbClient := db.GetMySQLDB()

	// Check if subscription with endpoint already exists (prevent 500 error on duplicate unique key index conflict)
	var existing db.WebPushSubscription
	if err := dbClient.Where("endpoint = ?", sub.Endpoint).First(&existing).Error; err == nil {
		log.Printf("WebPush subscription already registered for endpoint: %s", sub.Endpoint)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Subscription created successfully"})
	}

	if err := dbClient.Create(subScriptionModel).Error; err != nil {
		log.Printf("Error saving subscription to DB: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save subscription"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Subscription created successfully"})
}
