package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/r1i2t3/agni/pkg/db"
)

type ApplicationAuthRequest struct {
	ApplicationToken  string `json:"application_token"`
	ApplicationSecret string `json:"application_secret"`
}

func ApplicationAuth(c *fiber.Ctx) error {
	var req ApplicationAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	app, err := db.GetApplicationByTokenAndSecret(req.ApplicationToken, req.ApplicationSecret)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	c.Locals("app", app)
	return c.Next()
}
