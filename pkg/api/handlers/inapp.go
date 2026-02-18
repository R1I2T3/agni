// pkg/api/handlers/inapp.go
package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/r1i2t3/agni/pkg/db"
)

// GetInAppNotifications retrieves notifications for a user
// GET /api/inapp/notifications?user_id=<user_id>&unread_only=true&limit=50&offset=0
func GetInAppNotifications(c *fiber.Ctx) error {

	applicationID := c.Locals("application_id").(string)

	// Get query parameters
	userID := c.Query("user_id")
	unreadOnly := c.Query("unread_only", "false") == "true"
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}
	//print all for debugging
	println(userID, unreadOnly, limit, offset)
	mySQLDB := db.GetMySQLDB()
	query := mySQLDB.Where("application_id = ? AND recipient = ? AND channel = ?",
		applicationID, userID, "InApp")

	if unreadOnly {
		query = query.Where("read = ?", false)
	}

	// Get total count
	var total int64
	query.Model(&db.Notification{}).Count(&total)

	// Get notifications
	var notifications []db.Notification
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch notifications",
		})
	}

	return c.JSON(fiber.Map{
		"notifications": notifications,
		"total":         total,
		"limit":         limit,
		"offset":        offset,
		"has_more":      offset+len(notifications) < int(total),
	})
}

// MarkAsRead marks a single notification as read
// PUT /api/inapp/notifications/:id/read
func MarkNotificationAsRead(c *fiber.Ctx) error {
	applicationID := c.Locals("application_id").(string)
	notificationID := c.Params("id")

	if notificationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "notification_id is required",
		})
	}

	// Parse UUID
	id, err := uuid.Parse(notificationID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notification ID",
		})
	}

	// Update notification
	mySQLDB := db.GetMySQLDB()
	now := time.Now()
	result := mySQLDB.Model(&db.Notification{}).
		Where("id = ? AND application_id = ? AND channel = ?", id, applicationID, "InApp").
		Updates(map[string]interface{}{
			"read":    true,
			"read_at": now,
		})

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to mark notification as read",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Notification not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notification marked as read",
	})
}

// MarkAllAsRead marks all notifications for a user as read
// PUT /api/inapp/notifications/read-all
func MarkAllNotificationsAsRead(c *fiber.Ctx) error {
	applicationID := c.Locals("application_id").(string)

	var req struct {
		UserID string `json:"user_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}

	// Update all unread notifications for this user
	mySQLDB := db.GetMySQLDB()
	now := time.Now()
	result := mySQLDB.Model(&db.Notification{}).
		Where("application_id = ? AND recipient = ? AND channel = ? AND read = ?",
			applicationID, req.UserID, "InApp", false).
		Updates(map[string]interface{}{
			"read":    true,
			"read_at": now,
		})

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to mark notifications as read",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "All notifications marked as read",
		"count":   result.RowsAffected,
	})
}

// GetUnreadCount returns the count of unread notifications
// GET /api/inapp/notifications/unread-count?user_id=<user_id>
func GetUnreadCount(c *fiber.Ctx) error {
	applicationID := c.Locals("application_id").(string)
	userID := c.Query("user_id")

	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}

	// Count unread notifications
	mySQLDB := db.GetMySQLDB()
	var count int64
	err := mySQLDB.Model(&db.Notification{}).
		Where("application_id = ? AND recipient = ? AND channel = ? AND read = ?",
			applicationID, userID, "InApp", false).
		Count(&count).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get unread count",
		})
	}

	return c.JSON(fiber.Map{
		"user_id":      userID,
		"unread_count": count,
	})
}
