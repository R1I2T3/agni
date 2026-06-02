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

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	}
	var err error
	MySQLDB, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// SQLite doesn't need standard ping check via SQLDB
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
	return nil
}

func PingMySQL() error {
	if MySQLDB == nil {
		return fmt.Errorf("database client not initialized")
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

	sqlDB, err := MySQLDB.DB()
	if err != nil {
		result["status"] = "unhealthy"
		result["error"] = err.Error()
		return result
	}

	if err := sqlDB.Ping(); err != nil {
		result["status"] = "unhealthy"
		result["error"] = err.Error()
		return result
	}

	result["ping"] = true
	result["status"] = "healthy"
	stats := sqlDB.Stats()
	result["open_connections"] = stats.OpenConnections
	return result
}
