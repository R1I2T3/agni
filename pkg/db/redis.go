package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

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
