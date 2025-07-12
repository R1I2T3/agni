package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/r1i2t3/agni/pkg/notification"
	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

// QueuedNotification represents a notification that has been queued for delivery
type QueuedNotification struct {
	ID            string                           `json:"id"`
	ApplicationID string                           `json:"application_id"`
	QueueID       string                           `json:"queue_id"`
	Channel       notification.NotificationChannel `json:"channel"`
	Provider      string                           `json:"provider,omitempty"`
	Recipient     string                           `json:"recipient"`
	Subject       string                           `json:"subject,omitempty"`
	Message       string                           `json:"message"`
	TemplateID    string                           `json:"template_id,omitempty"`
	Status        string                           `json:"status"`
	Attempts      int                              `json:"attempts"`
	CreatedAt     time.Time                        `json:"created_at"`
	QueuedAt      time.Time                        `json:"queued_at"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host         string
	Port         string
	Password     string
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// InitRedis initializes Redis connection
func InitRedis(config RedisConfig) error {
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     config.Password,
		DB:           config.DB,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		// Connection pool settings
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test connection
	if err := PingRedis(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Printf("✅ Redis connected successfully at %s", addr)
	return nil
}

// PingRedis tests Redis connection
func PingRedis() error {
	pong, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	if pong != "PONG" {
		return fmt.Errorf("unexpected ping response: %s", pong)
	}

	return nil
}

// GetRedisClient returns the Redis client instance
func GetRedisClient() *redis.Client {
	return RedisClient
}

// CloseRedis closes Redis connection
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// RedisHealthCheck checks Redis health status
func RedisHealthCheck() map[string]interface{} {
	result := map[string]interface{}{
		"status":    "unknown",
		"connected": false,
		"ping":      false,
		"error":     nil,
	}

	if RedisClient == nil {
		result["status"] = "disconnected"
		result["error"] = "Redis client not initialized"
		return result
	}

	result["connected"] = true

	// Test ping
	pong, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		result["status"] = "unhealthy"
		result["error"] = err.Error()
		return result
	}

	if pong == "PONG" {
		result["ping"] = true
		result["status"] = "healthy"
	}

	// Get additional info
	info, err := RedisClient.Info(ctx, "server").Result()
	if err == nil {
		result["server_info"] = info
	}

	return result
}

func EnqueueNotification(Notification notification.Notification) (string, error) {

	QueueID := fmt.Sprintf("%s:%s:%s", Notification.ApplicationID, Notification.ID, Notification.Channel)
	queuedNotification := QueuedNotification{
		ID:            Notification.ID,
		ApplicationID: Notification.ApplicationID,
		QueueID:       QueueID,
		Channel:       Notification.Channel,
		Provider:      Notification.Provider,
		Recipient:     Notification.Recipient,
		Subject:       Notification.Subject,
		Message:       Notification.Message,
		TemplateID:    Notification.TemplateID,
		Status:        Notification.Status,
		Attempts:      Notification.Attempts,
		CreatedAt:     Notification.CreatedAt,
		QueuedAt:      time.Now(),
	}
	// Serialize the notification
	data, err := json.Marshal(queuedNotification)
	if err != nil {
		return "", fmt.Errorf("failed to serialize notification: %w", err)
	}
	// enqueue the notification in Redis
	err = RedisClient.LPush(ctx, "NotificationQueue", data).Err()
	if err != nil {
		return "", fmt.Errorf("failed to enqueue notification: %w", err)
	}
	log.Printf("✅ Notification queued successfully: %s", Notification.ID)
	return QueueID, nil
}

func DequeueNotification() (notification.Notification, error) {
	data, err := RedisClient.BRPop(ctx, 5*time.Second, "NotificationQueue").Result()
	if err != nil {
		return notification.Notification{}, fmt.Errorf("failed to dequeue notification: %w", err)
	}

	var notiNotification notification.Notification
	if err := json.Unmarshal([]byte(data[1]), &notiNotification); err != nil {
		return notification.Notification{}, fmt.Errorf("failed to deserialize notification: %w", err)
	}

	log.Printf("✅ Notification dequeued successfully: %s", notiNotification.ID)
	return notiNotification, nil
}
