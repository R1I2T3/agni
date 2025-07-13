package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/r1i2t3/agni/pkg/config"
	"github.com/r1i2t3/agni/pkg/db"
	"github.com/r1i2t3/agni/pkg/utils"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateTokenRequest struct {
	ApplicationName string `json:"application_name"`
}

type UpdateApplicationRequest struct {
	ApplicationName string `json:"application_name"`
}

func AdminLogin(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	fmt.Println("Received login request:", req)
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

func CreateApplicationAndApiTokenAndSecret(c *fiber.Ctx) error {
	var req CreateTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	fmt.Println("Received create application request:", req)
	// Logic to create application and generate API token and secret
	apiToken, apiSecret, err := utils.GenerateAppCredentials()
	if err != nil {
		return c.JSON(fiber.Map{"error": "Failed to generate credentials"})
	}
	err = db.CreateApplicationAndApiTokenAndSecret(req.ApplicationName, apiToken, apiSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create application"})
	}
	return c.JSON(fiber.Map{
		"api-token":  apiToken,
		"api-secret": apiSecret,
	})
}

func GetAllApplication(c *fiber.Ctx) error {
	applications, err := db.GetAllApplications()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch applications"})
	}
	return c.JSON(fiber.Map{"applications": applications})
}

func RegenerateToken(c *fiber.Ctx) error {
	var req UpdateApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	newToken, newSecret, err := utils.GenerateAppCredentials()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate new credentials"})
	}
	err = db.UpdateApplicationTokenAndSecret(req.ApplicationName, newToken, newSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update application credentials"})
	}
	return c.JSON(fiber.Map{"message": "Application credentials updated successfully", "api_token": newToken, "api_secret": newSecret})
}

func DeleteApplication(c *fiber.Ctx) error {
	var req UpdateApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	err := db.DeleteApplication(req.ApplicationName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete application"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Application deleted successfully"})
}
