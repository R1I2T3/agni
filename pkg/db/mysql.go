package db

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	MySQLDB *gorm.DB
)

type MySQLConfig struct {
	DSN      string
	LogLevel logger.LogLevel
}

func InitMySQL(config MySQLConfig, models ...interface{}) error {
	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	}
	var err error
	MySQLDB, err = gorm.Open(mysql.Open(config.DSN), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL database: %w", err)
	}
	if err := PingMySQL(); err != nil {
		return fmt.Errorf("failed to ping MySQL database: %w", err)
	}

	if len(models) > 0 {
		if err := MySQLDB.AutoMigrate(models...); err != nil {
			return fmt.Errorf("failed to run auto migrations: %w", err)
		}
		log.Printf("✅ Auto migrations completed successfully")
	}

	log.Printf("✅ MySQL connected successfully with DSN %s", config.DSN)
	return nil
}

func PingMySQL() error {
	if MySQLDB == nil {
		return fmt.Errorf("MySQL client not initialized")
	}

	sqlDB, err := MySQLDB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// GetMySQLDB returns the MySQL database instance
func GetMySQLDB() *gorm.DB {
	return MySQLDB
}

// CloseMySQL closes MySQL connection
func CloseMySQL() error {
	if MySQLDB != nil {
		sqlDB, err := MySQLDB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// MySQLHealthCheck checks MySQL health status
func MySQLHealthCheck() map[string]interface{} {
	result := map[string]interface{}{
		"status":    "unknown",
		"connected": false,
		"ping":      false,
		"error":     nil,
	}

	if MySQLDB == nil {
		result["status"] = "disconnected"
		result["error"] = "MySQL database not initialized"
		return result
	}

	result["connected"] = true

	// Test ping
	if err := PingMySQL(); err != nil {
		result["status"] = "unhealthy"
		result["error"] = err.Error()
		return result
	}

	result["ping"] = true
	result["status"] = "healthy"

	// Get database stats
	sqlDB, err := MySQLDB.DB()
	if err == nil {
		stats := sqlDB.Stats()
		result["open_connections"] = stats.OpenConnections
		result["in_use"] = stats.InUse
		result["idle"] = stats.Idle
	}

	return result
}
