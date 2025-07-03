package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors" // Add this import
	"github.com/r1i2t3/agni/pkg/db"
	"gorm.io/gorm/logger"
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

	// Initialize SQLite connection
	sqliteConfig := db.SQLiteConfig{
		DatabasePath: "./data/agni.db",
		LogLevel:     logger.Info, // or logger.Silent for production
	}

	if err := db.InitSQLite(sqliteConfig); err != nil {
		log.Fatalf("Failed to initialize SQLite: %v", err)
	}
	// Check health of both databases
	redisHealth := db.RedisHealthCheck()
	sqliteHealth := db.SQLiteHealthCheck()

	log.Printf("Redis Health Check: %+v", redisHealth)
	log.Printf("SQLite Health Check: %+v", sqliteHealth)

	// Check if both databases are healthy
	if redisHealth["ping"] == true {
		log.Println("✅ Redis is healthy")
	} else {
		log.Println("❌ Redis health check failed")
		log.Fatalf("Redis health check failed: %v", redisHealth)
	}

	if sqliteHealth["ping"] == true {
		log.Println("✅ SQLite is healthy")
	} else {
		log.Println("❌ SQLite health check failed")
		log.Fatalf("SQLite health check failed: %v", sqliteHealth)
	}

	app := fiber.New()

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*", // Allow all origins
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		MaxAge:           300, // Max age of preflight requests in seconds
	}))

	//temporary endpoint to test server is running
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Agni Notification Engine",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	//temporary endpoint to test Redis and SQLite connections for now
	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		redisHealth := db.RedisHealthCheck()
		sqliteHealth := db.SQLiteHealthCheck()

		overallStatus := "healthy"
		statusCode := 200

		if redisHealth["status"] != "healthy" || sqliteHealth["status"] != "healthy" {
			overallStatus = "unhealthy"
			statusCode = 503
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"service":   "Agni Notification Engine",
			"status":    overallStatus,
			"redis":     redisHealth,
			"sqlite":    sqliteHealth,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	log.Println("🚀 Starting Agni server on port 3000")
	log.Fatal(app.Listen(":3000"))
}
