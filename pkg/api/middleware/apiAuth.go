package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/r1i2t3/agni/pkg/db"
)

func APIKeyAuth(c *fiber.Ctx) error {
	apiKey := c.Get("X-API-Key")
	if apiKey == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "API key required",
		})
	}

	// Validate API key against applications table
	var app db.Application
	mySQLDB := db.GetMySQLDB()
	err := mySQLDB.Where("APIToken = ?", apiKey).First(&app).Error
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid API key",
		})
	}

	// Store application in context
	c.Locals("application", &app)
	return c.Next()
}
