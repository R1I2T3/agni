// pkg/inapp/middleware/websocket-auth.go
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/r1i2t3/agni/pkg/utils"
)

// WebSocketJWTAuth validates JWT token for WebSocket connections
// Accepts token from either:
// - Query parameter: ?token=<jwt>
// - Cookie: app_token
func WebSocketJWTAuth(c *fiber.Ctx) error {
	// Try to get token from query parameter first
	token := c.Query("token")

	// If not in query, try to get from cookie
	if token == "" {
		token = c.Cookies("Agni-auth-token")
	}

	// If still no token found
	if token == "" {
		return websocket.New(func(conn *websocket.Conn) {
			_ = conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure,
					"Authentication required: provide token via query parameter or cookie"))
			conn.Close()
		})(c)
	}

	// Validate and extract claims from JWT
	applicationID, userID, err := utils.ValidateApplicationJWT(token)
	if err != nil {
		return websocket.New(func(conn *websocket.Conn) {
			_ = conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure,
					"Invalid token: "+err.Error()))
			conn.Close()
		})(c)
	}

	// Store in context for handler
	c.Locals("application_id", applicationID.String())
	c.Locals("user_id", userID)

	return c.Next()
}
