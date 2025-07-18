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

	// Health check endpoints
	app.Get("/health", handlers.HealthCheck)

	// admin auth routes
	app.Post("/api/admin/auth/login", handlers.AdminLogin)
	app.Post("/api/admin/logout", handlers.AdminLogout)
	app.Get("/api/admin/dashboard", middleware.RequireAdmin, handlers.AdminDashBoardRedirect)

	// admin functions routes
	app.Post("/api/admin/create-application", middleware.RequireAdmin, handlers.CreateApplicationAndApiTokenAndSecret)
	app.Get("/api/admin/applications", middleware.RequireAdmin, handlers.GetAllApplication)
	app.Put("/api/admin/regenerate-token", middleware.RequireAdmin, handlers.RegenerateToken)
	app.Delete("/api/admin/delete-application", middleware.RequireAdmin, handlers.DeleteApplication)

	// notification routes
	app.Post("/api/notification/send", middleware.ApplicationAuth, handlers.EnqueueNotification)
}
