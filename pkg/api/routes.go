package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/r1i2t3/agni/pkg/api/handlers"
	"github.com/r1i2t3/agni/pkg/api/middleware"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Agni Notification Engine",
			"version": "1.0.0",
			"status":  "running",
		})
	})
	// health check end points
	app.Get("/health", handlers.HealthCheck)

	// admin auth routes
	app.Post("/admin/auth/login", handlers.AdminLogin)
	app.Post("/admin/logout", handlers.AdminLogout)
	app.Get("/admin/dashboard", middleware.RequireAdmin, handlers.AdminDashBoardRedirect)
}
