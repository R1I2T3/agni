package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/r1i2t3/agni/pkg/config"
	"github.com/r1i2t3/agni/pkg/db"
	inapp "github.com/r1i2t3/agni/pkg/inapp"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load environment configurations
	envConfig := config.GetEnvConfig()
	redisConfig := db.RedisConfig{
		Host:         envConfig.RedisEnvConfig.Host,
		Port:         envConfig.RedisEnvConfig.Port,
		Password:     envConfig.RedisEnvConfig.Password,
		DB:           envConfig.RedisEnvConfig.DB,
		DialTimeout:  envConfig.RedisEnvConfig.DialTimeout,
		ReadTimeout:  envConfig.RedisEnvConfig.ReadTimeout,
		WriteTimeout: envConfig.RedisEnvConfig.WriteTimeout,
	}

	// Initialize databases
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

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("ok")
	})

	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		user := conn.Query("user")
		if user == "" {
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "user required"))
			return
		}

		client := inapp.DefaultHub.Register(user, conn)
		client.ReadPump() // Block handler to keep connection open
	}))

	addr := ":" + envConfig.InAppServiceConfig.Port
	if addr == ":" {
		addr = ":4000"
	}
	log.Printf("inapp: listening on %s", addr)
	log.Fatal(app.Listen(addr))
}

// startBroadcastSubscriber listens to Redis Pub/Sub and delivers to local WebSocket hub
func startBroadcastSubscriber(ctx context.Context, rdb *redis.Client) {
	// Pattern subscribe to all broadcast channels: inapp:broadcast:*
	pubsub := rdb.PSubscribe(ctx, inapp.BroadcastChannelPrefix+"*")
	defer pubsub.Close()

	log.Println("📡 Broadcast subscriber started, listening for notifications...")

	ch := pubsub.Channel()
	for msg := range ch {
		// msg.Channel = "inapp:broadcast:user123"
		// msg.Payload = JSON notification

		// Extract recipient from channel name
		parts := strings.Split(msg.Channel, ":")
		if len(parts) < 3 {
			log.Printf("invalid broadcast channel format: %s", msg.Channel)
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

	log.Println("⚠️  Broadcast subscriber stopped")
}
