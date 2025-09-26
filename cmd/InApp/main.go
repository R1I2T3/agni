package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/r1i2t3/agni/pkg/config"
	"github.com/r1i2t3/agni/pkg/db"
	inapp "github.com/r1i2t3/agni/pkg/inapp"
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

	// initial Cunsumers
	ctx := context.Background()

	// Create consumer group if not exists
	rdb := db.GetRedisClient()

	_ = rdb.XGroupCreateMkStream(ctx, inapp.StreamName, inapp.GroupName, "$").Err()

	// start consumer loop
	go inapp.StartConsumer(ctx, rdb, inapp.GroupName, "inapp-consumer-1")

	// Fiber HTTP + WebSocket server
	app := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("ok")
	})

	// Replace adaptor route with fiber websocket route
	// app.Get("/ws", adaptor.HTTPHandler(http.HandlerFunc(inapp.WSHandler)))
	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		// extract user query param
		user := string(conn.Params("user")) // fallback; prefer Query if available

		// register connection with your hub: pass the fiber websocket.Conn directly
		client := inapp.DefaultHub.Register(user, conn)
		_ = client
		client.ReadPump()
	}))

	addr := ":" + os.Getenv("INAPP_PORT")
	if addr == ":" {
		addr = ":4000"
	}
	log.Printf("inapp: listening on %s", addr)
	log.Fatal(app.Listen(addr))

}
