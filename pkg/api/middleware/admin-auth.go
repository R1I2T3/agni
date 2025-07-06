package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/r1i2t3/agni/pkg/utils"
)

func RequireAdmin(c *fiber.Ctx) error {
	cookie := c.Cookies("admin_token")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	token, claims, err := utils.ParseJWT(cookie)
	if err != nil || !token.Valid || claims["admin"] != true {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	c.Locals("admin", claims["username"])
	return c.Next()
}
