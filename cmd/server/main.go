package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/r1i2t3/agni/pkg/db"
)

func main() {

	// Initialize Redis connection
	redisConfig := db.RedisConfig{
		Host:         "redis",
		Port:         "6379",
		Password:     "",
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	if err := db.InitRedis(redisConfig); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	RedisHealthInfo := db.RedisHealthCheck()

	log.Printf("Redis Health Check: %+v", RedisHealthInfo)

	// Check if Redis is healthy
	if RedisHealthInfo["ping"] == true {
		log.Println("✅ Redis is healthy")
	} else {
		log.Println("❌ Redis health check failed or Redis is not healthy")
		log.Fatalf("Redis health check failed: %v", RedisHealthInfo)

		//terminate the application if Redis is not healthy
		return
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server is running!")
	})

	log.Fatal(app.Listen(":3000"))
}
