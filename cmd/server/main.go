package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors" // Add this import
	"github.com/r1i2t3/agni/pkg/config"           // Import your config package if needed
	"github.com/r1i2t3/agni/pkg/db"
)

// reformat  if looks messy
func main() {

	//get environment variables for Redis and SQLite configurations and server configurations
	EnvConfig := config.GetEnvConfig()
	RedisConfig := EnvConfig.RedisEnvConfig
	SQLiteConfig := EnvConfig.SQLiteEnvConfig
	ServerConfig := EnvConfig.ServerEnvConfig
	CorsConfig := EnvConfig.CorsEnvConfig

	// Initialize Redis connection
	redisConfig := db.RedisConfig{
		Host:         RedisConfig.Host,
		Port:         RedisConfig.Port,
		Password:     RedisConfig.Password,
		DB:           RedisConfig.DB,
		DialTimeout:  RedisConfig.DialTimeout,
		ReadTimeout:  RedisConfig.ReadTimeout,
		WriteTimeout: RedisConfig.WriteTimeout,
	}

	if err := db.InitRedis(redisConfig); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// Initialize SQLite connection
	sqliteConfig := db.SQLiteConfig{
		DatabasePath: SQLiteConfig.DatabasePath,
		LogLevel:     config.GetLogLevel(SQLiteConfig.LogLevel), // default to "info" if not set
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
		AllowOrigins: CorsConfig.AllowOrigins,
		AllowMethods: CorsConfig.AllowMethods,
		AllowHeaders: CorsConfig.AllowHeaders,
		MaxAge:       CorsConfig.MaxAge, // Max age of preflight requests in seconds
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

	log.Println("🚀 Starting Agni server on port ", ServerConfig.Port)
	log.Fatal(app.Listen(":" + ServerConfig.Port))
}
