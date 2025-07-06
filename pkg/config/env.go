package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm/logger"
)

func GetEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func GetEnvAsInt(key string, defaultValue int) int {
	valueStr := GetEnv(key, strconv.Itoa(defaultValue))
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func GetEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := GetEnv(key, fmt.Sprintf("%d", defaultValue))
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
func GetEnvAsBool(key string, defaultValue bool) bool {
	valueStr := GetEnv(key, fmt.Sprintf("%t", defaultValue))
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

type RedisEnvConfig struct {
	Host         string
	Port         string
	Password     string
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func GetRedisEnvConfig() RedisEnvConfig {
	return RedisEnvConfig{
		Host:         GetEnv("REDIS_HOST", "localhost"),
		Port:         GetEnv("REDIS_PORT", "6379"),
		Password:     GetEnv("REDIS_PASSWORD", ""),
		DB:           GetEnvAsInt("REDIS_DB", 0),
		DialTimeout:  GetEnvAsDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
		ReadTimeout:  GetEnvAsDuration("REDIS_READ_TIMEOUT", 5*time.Second),
		WriteTimeout: GetEnvAsDuration("REDIS_WRITE_TIMEOUT", 5*time.Second),
	}
}

type SQLiteEnvConfig struct {
	DatabasePath string
	LogLevel     string
}

func GetSQLiteEnvConfig() SQLiteEnvConfig {
	return SQLiteEnvConfig{
		DatabasePath: GetEnv("SQLITE_DATABASE_PATH", "./data/agni.db"),
		LogLevel:     GetEnv("SQLITE_LOG_LEVEL", "info"), // Default to "info" if not set
	}
}

type ServerEnvConfig struct {
	Port string
}

type AdminEnvConfig struct {
	Admin_Username string
	Admin_Password string
	JWT_Secret     string
}

func GetServerEnvConfig() ServerEnvConfig {
	return ServerEnvConfig{
		Port: GetEnv("SERVER_PORT", "3000"),
	}
}

func GetAdminEnvConfig() AdminEnvConfig {
	return AdminEnvConfig{
		Admin_Username: GetEnv("ADMIN_USERNAME", "admin"),
		Admin_Password: GetEnv("ADMIN_PASSWORD", "admin123"),
		JWT_Secret:     GetEnv("JWT_SECRET", ""),
	}
}

type CorsEnvConfig struct {
	AllowOrigins     string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
	MaxAge           int
}

func GetCorsEnvConfig() CorsEnvConfig {
	return CorsEnvConfig{
		AllowOrigins:     GetEnv("CORS_ORIGINS", "*"),
		AllowMethods:     GetEnv("CORS_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS"),
		AllowHeaders:     GetEnv("CORS_HEADERS", "Content-Type,Authorization"),
		AllowCredentials: GetEnvAsBool("CORS_CREDENTIALS", false),
		MaxAge:           GetEnvAsInt("CORS_MAX_AGE", 300), // Default to 300 seconds
	}
}

type EnvConfig struct {
	ServerEnvConfig ServerEnvConfig
	SQLiteEnvConfig SQLiteEnvConfig
	RedisEnvConfig  RedisEnvConfig
	CorsEnvConfig   CorsEnvConfig
	AdminEnvConfig  AdminEnvConfig
}

func GetEnvConfig() EnvConfig {
	return EnvConfig{
		ServerEnvConfig: GetServerEnvConfig(),
		SQLiteEnvConfig: GetSQLiteEnvConfig(),
		RedisEnvConfig:  GetRedisEnvConfig(),
		CorsEnvConfig:   GetCorsEnvConfig(),
		AdminEnvConfig:  GetAdminEnvConfig(),
	}
}

func GetLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info // Default to Info if unknown
	}
}
