package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm/logger"
)

type EnvConfig struct {
	ServerEnvConfig  ServerEnvConfig
	RedisEnvConfig   RedisEnvConfig
	CorsEnvConfig    CorsEnvConfig
	AdminEnvConfig   AdminEnvConfig
	EmailEnvConfig   EmailEnvConfig
	ResendEnvConfig  ResendEnvConfig
	MySQLConfig      MySQLConfig
	TwilioEnvConfig  TwilioEnvConfig
	WebPushEnvConfig WebPushEnvConfig
	InAppConfig      InAppConfig
}

func GetEnvConfig() EnvConfig {
	return EnvConfig{
		ServerEnvConfig:  GetServerEnvConfig(),
		MySQLConfig:      GetMySQLDBConfig(),
		RedisEnvConfig:   GetRedisEnvConfig(),
		CorsEnvConfig:    GetCorsEnvConfig(),
		AdminEnvConfig:   GetAdminEnvConfig(),
		EmailEnvConfig:   GetEmailEnvConfig(),
		ResendEnvConfig:  GetResendEnvConfig(),
		TwilioEnvConfig:  GetTwilioEnvConfig(),
		WebPushEnvConfig: GetWebPushEnvConfig(),
		InAppConfig:      GetInAppConfig(),
	}
}

type ServerEnvConfig struct {
	Port string
}

func GetServerEnvConfig() ServerEnvConfig {
	return ServerEnvConfig{
		Port: GetEnv("SERVER_PORT", "8080"),
	}
}

type MySQLConfig struct {
	MYSQL_USER          string
	MYSQL_ROOT_PASSWORD string
	DB_HOST             string
	MYSQL_DATABASE      string
}

func GetMySQLDBConfig() MySQLConfig {
	return MySQLConfig{
		MYSQL_USER:          GetEnv("MYSQL_USER", ""),
		MYSQL_ROOT_PASSWORD: GetEnv("MYSQL_ROOT_PASSWORD", ""),
		DB_HOST:             GetEnv("DB_HOST", ""),
		MYSQL_DATABASE:      GetEnv("MYSQL_DATABASE", ""),
	}
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

type AdminEnvConfig struct {
	Admin_Username string
	Admin_Password string
	JWT_Secret     string
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

func GetAdminEnvConfig() AdminEnvConfig {
	return AdminEnvConfig{
		Admin_Username: GetEnv("ADMIN_USERNAME", "admin"),
		Admin_Password: GetEnv("ADMIN_PASSWORD", "admin123"),
		JWT_Secret:     GetEnv("JWT_SECRET", ""),
	}
}

type EmailEnvConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromAddress  string
	UseTLS       bool
}

func GetEmailEnvConfig() EmailEnvConfig {
	return EmailEnvConfig{
		SMTPHost:     GetEnv("EMAIL_SMTP_HOST", "smtp.example.com"),
		SMTPPort:     GetEnv("EMAIL_SMTP_PORT", "587"),
		SMTPUsername: GetEnv("EMAIL_SMTP_USERNAME", ""),
		SMTPPassword: GetEnv("EMAIL_SMTP_PASSWORD", ""),
		FromAddress:  GetEnv("EMAIL_FROM_ADDRESS", "noreply@example.com"),
		UseTLS:       GetEnvAsBool("EMAIL_USE_TLS", true),
	}
}

type ResendEnvConfig struct {
	APIKey      string
	FromAddress string
}

func GetResendEnvConfig() ResendEnvConfig {
	return ResendEnvConfig{
		APIKey:      GetEnv("RESEND_API_KEY", ""),
		FromAddress: GetEnv("RESEND_FROM_ADDRESS", ""),
	}
}

type TwilioEnvConfig struct {
	TWILIO_PHONE_NUMBER string
	ACCOUNT_SID         string
	AUTH_TOKEN          string
}

func GetTwilioEnvConfig() TwilioEnvConfig {
	return TwilioEnvConfig{
		TWILIO_PHONE_NUMBER: GetEnv("TWILIO_PHONE_NUMBER", ""),
		ACCOUNT_SID:         GetEnv("TWILIO_ACCOUNT_SID", ""),
		AUTH_TOKEN:          GetEnv("TWILIO_AUTH_TOKEN", ""),
	}
}

type WebPushEnvConfig struct {
	VAPID_PUBLIC_KEY  string
	VAPID_PRIVATE_KEY string
	VAPID_SUBJECT     string
}

func GetWebPushEnvConfig() WebPushEnvConfig {
	return WebPushEnvConfig{
		VAPID_PUBLIC_KEY:  GetEnv("VAPID_PUBLIC_KEY", ""),
		VAPID_PRIVATE_KEY: GetEnv("VAPID_PRIVATE_KEY", ""),
		VAPID_SUBJECT:     GetEnv("VAPID_SUBJECT", ""),
	}
}

type InAppConfig struct {
	stream string
}

func GetInAppConfig() InAppConfig {
	return InAppConfig{
		stream: GetEnv("Redis_InApp_streamName", "stream:inapp")}
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
