package main

import (
	"context"
	"log"
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

	_ = rdb.XGroupCreateMkStream(ctx, envConfig.InAppServiceConfig.StreamName, envConfig.InAppServiceConfig.GroupName, "$").Err()

	// start consumer loop
	go inapp.StartConsumer(ctx, rdb, envConfig.InAppServiceConfig.GroupName, envConfig.InAppServiceConfig.ConsumerName)

	// Fiber HTTP + WebSocket server
	app := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("ok")
	})

	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {

		user := string(conn.Query("user"))

		client := inapp.DefaultHub.Register(user, conn)
		_ = client
		client.ReadPump()
	}))

	addr := ":" + envConfig.InAppServiceConfig.Port
	if addr == ":" {
		addr = ":4000"
	}
	log.Printf("inapp: listening on %s", addr)
	log.Fatal(app.Listen(addr))

}
