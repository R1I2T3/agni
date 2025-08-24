package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/r1i2t3/agni/pkg/api"
	"github.com/r1i2t3/agni/pkg/config"
	"github.com/r1i2t3/agni/pkg/db"
	workers "github.com/r1i2t3/agni/pkg/queue/Workers"
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
	fmt.Println("Redis Config:", redisConfig)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_ROOT_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("MYSQL_DATABASE"),
	)
	fmt.Println("MySQL DSN:", dsn)
	mySQLConfig := db.MySQLConfig{
		DSN: dsn,
	}

	// Initialize databases
	config.InitializeRedis(redisConfig)
	config.InitializeMySQL(mySQLConfig)

	// Initialize channels
	config.InitializeEmailChannel(&envConfig.EmailEnvConfig)
	config.InitializeResendProvider(&envConfig.ResendEnvConfig)

	// Create Fiber app
	app := fiber.New()

	// Configure CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: envConfig.CorsEnvConfig.AllowOrigins,
		AllowMethods: envConfig.CorsEnvConfig.AllowMethods,
		AllowHeaders: envConfig.CorsEnvConfig.AllowHeaders,
		MaxAge:       envConfig.CorsEnvConfig.MaxAge,
	}))

	// Setup routes
	api.SetupRoutes(app)

	// start of Notification Workers
	workerPool := workers.NewWorkerPool(5, "QueuedNotification")
	workerPool.Start()

	// Start Delayed Queue Processor
	log.Println("⏰ Starting delayed queue processor...")
	delayedProcessor := workers.NewDelayedQueueProcessor(
		"QueuedNotification:delayed", // Delayed queue name
		"QueuedNotification",         // Main queue name
		time.Second*10,               // Check every 10 seconds
	)
	delayedProcessor.Start()

	// Start server
	log.Println("🚀 Starting Agni server on port", envConfig.ServerEnvConfig.Port)
	log.Fatal(app.Listen(":" + envConfig.ServerEnvConfig.Port))
}
