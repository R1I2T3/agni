package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/r1i2t3/agni/pkg/config"
	"github.com/r1i2t3/agni/pkg/db"
	inapp "github.com/r1i2t3/agni/pkg/inapp"
	"github.com/r1i2t3/agni/pkg/inapp/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load environment configurations
	envConfig := config.GetEnvConfig()

	// Initialize Redis
	redisConfig := db.RedisConfig{
		Host:         envConfig.RedisEnvConfig.Host,
		Port:         envConfig.RedisEnvConfig.Port,
		Password:     envConfig.RedisEnvConfig.Password,
		DB:           envConfig.RedisEnvConfig.DB,
		DialTimeout:  envConfig.RedisEnvConfig.DialTimeout,
		ReadTimeout:  envConfig.RedisEnvConfig.ReadTimeout,
		WriteTimeout: envConfig.RedisEnvConfig.WriteTimeout,
	}
	config.InitializeRedis(redisConfig)

	ctx := context.Background()
	rdb := db.GetRedisClient()

	// Create consumer group if not exists
	_ = rdb.XGroupCreateMkStream(ctx, envConfig.InAppServiceConfig.StreamName, envConfig.InAppServiceConfig.GroupName, "$").Err()

	// Start consumer loop (reads from stream, publishes to pub/sub)
	go inapp.StartConsumer(ctx, rdb, envConfig.InAppServiceConfig.GroupName, envConfig.InAppServiceConfig.ConsumerName)

	// Start broadcast subscriber (listens to pub/sub, delivers to local WebSocket hub)
	go startBroadcastSubscriber(ctx, rdb)

	// Fiber HTTP + WebSocket server
	app := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	})

	// Middleware
	app.Use(cors.New())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("ok")
	})

	// WebSocket endpoint - protected with JWT
	app.Get("/ws", middleware.WebSocketJWTAuth, websocket.New(func(conn *websocket.Conn) {
		userID := conn.Locals("user_id").(string)

		client := inapp.DefaultHub.Register(userID, conn)
		client.ReadPump()
	}))

	addr := ":" + envConfig.InAppServiceConfig.Port
	if addr == ":" {
		addr = ":4000"
	}
	log.Printf("🚀 InApp WebSocket service listening on %s", addr)
	log.Printf("   - GET /ws?token=<jwt> (WebSocket with JWT authentication)")
	log.Fatal(app.Listen(addr))
}

// startBroadcastSubscriber listens to Redis Pub/Sub and delivers to local WebSocket hub
func startBroadcastSubscriber(ctx context.Context, rdb *redis.Client) {
	// Pattern subscribe to all broadcast channels: inapp:broadcast:*
	pubsub := rdb.PSubscribe(ctx, inapp.BroadcastChannelPrefix+"*")
	defer pubsub.Close()

	log.Println("📡 Broadcast subscriber started")

	ch := pubsub.Channel()
	for msg := range ch {
		parts := strings.Split(msg.Channel, ":")
		if len(parts) < 3 {
			log.Printf("invalid broadcast channel: %s", msg.Channel)
			continue
		}
		recipient := parts[2]

		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
			log.Printf("invalid broadcast payload: %v", err)
			continue
		}

		// Check if this container has the user's WebSocket connection
		// BroadcastToUser will check the local hub and only send if connected here
		inapp.DefaultHub.BroadcastToUser(recipient, payload)
	}
}
