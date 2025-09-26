package handlers

import (
	"fmt"

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
	fmt.Println("Handling web push subscription")
	if err := c.BodyParser(&sub); err != nil {
		fmt.Println("Error parsing body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	fmt.Println("Received subscription:", sub)
	var subScriptionModel = &db.WebPushSubscription{
		Endpoint: sub.Endpoint,
		Auth:     sub.Keys.Auth,
		P256dh:   sub.Keys.P256dh,
		Device:   sub.Device,
		UserID:   sub.UserID,
	}
	db := db.GetMySQLDB()
	if err := db.Create(&subScriptionModel).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save subscription"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Subscription created successfully"})
}
