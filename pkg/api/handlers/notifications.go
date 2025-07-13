package handlers

import "github.com/gofiber/fiber/v2"

func SendNotification(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Notification sent successfully",
		"status":  "success",
	})
}
