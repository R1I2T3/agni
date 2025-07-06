package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/r1i2t3/agni/pkg/api/handlers"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Agni Notification Engine",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	app.Get("/health", handlers.HealthCheck)
}
