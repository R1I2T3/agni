package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	SQLiteDB *gorm.DB
)

// SQLiteConfig holds SQLite configuration
type SQLiteConfig struct {
	DatabasePath string
	LogLevel     logger.LogLevel
}

// InitSQLite initializes SQLite connection
func InitSQLite(config SQLiteConfig, models ...interface{}) error {
	// Ensure directory exists
	if err := ensureDirectoryExists(config.DatabasePath); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	}

	var err error
	SQLiteDB, err = gorm.Open(sqlite.Open(config.DatabasePath), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	// Test connection
	if err := PingSQLite(); err != nil {
		return fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	// Run auto migrations
	if len(models) > 0 {
		if err := SQLiteDB.AutoMigrate(models...); err != nil {
			return fmt.Errorf("failed to run auto migrations: %w", err)
		}
		log.Printf("✅ Auto migrations completed successfully")
	}

	log.Printf("✅ SQLite connected successfully at %s", config.DatabasePath)
	return nil
}

// PingSQLite tests SQLite connection
func PingSQLite() error {
	if SQLiteDB == nil {
		return fmt.Errorf("SQLite client not initialized")
	}

	sqlDB, err := SQLiteDB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// GetSQLiteDB returns the SQLite database instance
func GetSQLiteDB() *gorm.DB {
	return SQLiteDB
}

// CloseSQLite closes SQLite connection
func CloseSQLite() error {
	if SQLiteDB != nil {
		sqlDB, err := SQLiteDB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// SQLiteHealthCheck checks SQLite health status
func SQLiteHealthCheck() map[string]interface{} {
	result := map[string]interface{}{
		"status":    "unknown",
		"connected": false,
		"ping":      false,
		"error":     nil,
	}

	if SQLiteDB == nil {
		result["status"] = "disconnected"
		result["error"] = "SQLite database not initialized"
		return result
	}

	result["connected"] = true

	// Test ping
	if err := PingSQLite(); err != nil {
		result["status"] = "unhealthy"
		result["error"] = err.Error()
		return result
	}

	result["ping"] = true
	result["status"] = "healthy"

	// Get database file info
	sqlDB, err := SQLiteDB.DB()
	if err == nil {
		stats := sqlDB.Stats()
		result["open_connections"] = stats.OpenConnections
		result["in_use"] = stats.InUse
		result["idle"] = stats.Idle
	}

	return result
}

// ensureDirectoryExists creates directory if it doesn't exist
func ensureDirectoryExists(filePath string) error {
	dir := getDirectory(filePath)
	if dir == "" {
		return nil // No directory to create
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// getDirectory extracts directory from file path
func getDirectory(filePath string) string {
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '/' || filePath[i] == '\\' {
			return filePath[:i]
		}
	}
	return ""
}
