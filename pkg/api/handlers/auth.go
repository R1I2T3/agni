// pkg/api/handlers/auth.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/r1i2t3/agni/pkg/db"
	"github.com/r1i2t3/agni/pkg/utils"
)

type HMACAuthRequest struct {
	APIToken  string `json:"api_token"`
	UserID    string `json:"user_id"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

type HMACAuthResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
	UserID    string `json:"user_id"`
}

// HMACLogin authenticates client with HMAC and returns JWT
// POST /api/auth/login
// Body: { "api_token": "...", "user_id": "...", "timestamp": 123456789, "signature": "..." }
// Signature = HMAC-SHA256(api_secret, "user_id:timestamp")
// timestamps is used to prevent replay attacks
func HMACLogin(c *fiber.Ctx) error {
	var req HMACAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.APIToken == "" || req.UserID == "" || req.Signature == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields: api_token, user_id, timestamp, signature",
		})
	}

	// Lookup application by API token
	var app db.Application
	mySQLDB := db.GetMySQLDB()
	if err := mySQLDB.Where("api_token = ?", req.APIToken).First(&app).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Validate HMAC signature with 5 minute expiry
	if err := utils.ValidateTimestampedHMAC(app.APISecret, req.UserID, req.Signature, req.Timestamp, 300); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication failed: " + err.Error(),
		})
	}

	// Generate JWT token
	token, err := utils.GenerateApplicationJWT(app.ID, req.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Set JWT as HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "Agni-auth-token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400,
		Secure:   false,
		HTTPOnly: true,  // Cannot be accessed by JavaScript
		SameSite: "Lax", // CSRF protection
	})

	// Return success response (no token in body)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":    true,
		"user_id":    req.UserID,
		"expires_in": 86400,
		"message":    "Authentication successful",
	})
}
