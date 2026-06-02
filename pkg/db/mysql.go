package db

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
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

func InitMySQL(mode string, config MySQLConfig, models ...interface{}) error {
	var dialector gorm.Dialector
	if mode == "local" {
		log.Println("🔌 GORM Mode: LOCAL (Initializing SQLite connection: agni.db)")
		dialector = sqlite.Open("agni.db")
	} else {
		log.Printf("🔌 GORM Mode: PROD (Initializing MySQL connection: %s)", config.DSN)
		dialector = mysql.Open(config.DSN)
	}

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	}
	var err error
	MySQLDB, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// SQLite doesn't need TCP socket-based ping check inside InitMySQL
	if mode != "local" {
		if err := PingMySQL(); err != nil {
			return fmt.Errorf("failed to ping MySQL database: %w", err)
		}
	}

	if len(models) > 0 {
		if err := MySQLDB.AutoMigrate(models...); err != nil {
			return fmt.Errorf("failed to run auto migrations: %w", err)
		}
		log.Printf("✅ Auto migrations completed successfully")
	}

	if mode == "local" {
		log.Println("✅ SQLite connected successfully")
	} else {
		log.Printf("✅ MySQL connected successfully with DSN %s", config.DSN)
	}
	return nil
}

func PingMySQL() error {
	if MySQLDB == nil {
		return fmt.Errorf("database client not initialized")
	}

	if MySQLDB.Dialector.Name() == "sqlite" {
		return nil
	}

	sqlDB, err := MySQLDB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// GetMySQLDB returns the database instance
func GetMySQLDB() *gorm.DB {
	return MySQLDB
}

// CloseMySQL closes database connection
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

// MySQLHealthCheck checks database health status
func MySQLHealthCheck() map[string]interface{} {
	result := map[string]interface{}{
		"status":    "unknown",
		"connected": false,
		"ping":      false,
		"error":     nil,
	}

	if MySQLDB == nil {
		result["status"] = "disconnected"
		result["error"] = "Database not initialized"
		return result
	}

	result["connected"] = true

	// Skip standard TCP socket Ping for SQLite
	if MySQLDB.Dialector.Name() == "sqlite" {
		result["ping"] = true
		result["status"] = "healthy"
		sqlDB, err := MySQLDB.DB()
		if err == nil {
			stats := sqlDB.Stats()
			result["open_connections"] = stats.OpenConnections
			result["in_use"] = stats.InUse
			result["idle"] = stats.Idle
		}
		return result
	}

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
