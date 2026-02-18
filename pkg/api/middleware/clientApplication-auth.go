// pkg/inapp/middleware/websocket-auth.go
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/r1i2t3/agni/pkg/utils"
)

func ClientApplicationAuth(c *fiber.Ctx) error {
	// Try to get token from query parameter first
	token := c.Query("token")

	// If not in query, try to get from cookie
	if token == "" {
		token = c.Cookies("Agni-auth-token")
	}

	// If still no token found
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required: provide token via query parameter or cookie",
		})
	}

	// Validate and extract claims from JWT
	applicationID, userID, err := utils.ValidateApplicationJWT(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token: " + err.Error(),
		})
	}

	// Store in context for handler
	c.Locals("application_id", applicationID.String())
	c.Locals("user_id", userID)

	return c.Next()
}
