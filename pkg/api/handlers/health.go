package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/r1i2t3/agni/pkg/db"
)

func HealthCheck(c *fiber.Ctx) error {
	redisHealth := db.RedisHealthCheck()
	mySQLHealth := db.MySQLHealthCheck()

	overallStatus := "healthy"
	statusCode := 200

	if redisHealth["status"] != "healthy" || mySQLHealth["status"] != "healthy" {
		overallStatus = "unhealthy"
		statusCode = 503
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"service":   "Agni Notification Engine",
		"status":    overallStatus,
		"redis":     redisHealth,
		"mysql":     mySQLHealth,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
