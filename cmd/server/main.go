package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/r1i2t3/agni/pkg/api"
	"github.com/r1i2t3/agni/pkg/config"
	"github.com/r1i2t3/agni/pkg/db"
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
	sqliteConfig := db.SQLiteConfig{
		DatabasePath: envConfig.SQLiteEnvConfig.DatabasePath,
		LogLevel:     config.GetLogLevel(envConfig.SQLiteEnvConfig.LogLevel),
	}

	// Initialize databases
	config.InitializeRedis(redisConfig)
	config.InitializeSQLite(sqliteConfig)

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

	// Start server
	log.Println("🚀 Starting Agni server on port", envConfig.ServerEnvConfig.Port)
	log.Fatal(app.Listen(":" + envConfig.ServerEnvConfig.Port))
}
