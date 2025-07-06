package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/r1i2t3/agni/pkg/config"
	"github.com/r1i2t3/agni/pkg/utils"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AdminLogin(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	envConfig := config.GetEnvConfig()
	if req.Username == envConfig.AdminEnvConfig.Admin_Username && req.Password == envConfig.AdminEnvConfig.Admin_Password {
		token, err := utils.GenerateAdminJWT(req.Username)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Token generation failed"})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "admin_token",
			Value:    token,
			HTTPOnly: true,
			Secure:   false,
			MaxAge:   60 * 60 * 24,
		})

		return c.JSON(fiber.Map{"message": "Login successful"})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
}

func AdminLogout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "admin_token",
		Value:    "",
		HTTPOnly: true,
		MaxAge:   -1,
	})
	return c.JSON(fiber.Map{"message": "Logged out"})
}

func AdminDashBoardRedirect(c *fiber.Ctx) error {
	user := c.Locals("admin")
	return c.JSON(fiber.Map{
		"message": "Welcome, admin",
		"user":    user,
	})

}
