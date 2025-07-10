package config

import (
	"log"

	"github.com/r1i2t3/agni/pkg/db"
)

func InitializeRedis(redisConfig db.RedisConfig) {
	if err := db.InitRedis(redisConfig); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	redisHealth := db.RedisHealthCheck()
	log.Printf("Redis Health Check: %+v", redisHealth)

	if redisHealth["ping"] == true {
		log.Println("✅ Redis is healthy")
	} else {
		log.Fatalf("❌ Redis health check failed: %v", redisHealth)
	}
}

func InitializeSQLite(sqliteConfig db.SQLiteConfig) {
	allModel := []interface{}{
		&db.Notification{},
		&db.Application{},
	}
	if err := db.InitSQLite(sqliteConfig, allModel...); err != nil {
		log.Fatalf("Failed to initialize SQLite: %v", err)
	}

	sqliteHealth := db.SQLiteHealthCheck()
	log.Printf("SQLite Health Check: %+v", sqliteHealth)

	if sqliteHealth["ping"] == true {
		log.Println("✅ SQLite is healthy")
	} else {
		log.Fatalf("❌ SQLite health check failed: %v", sqliteHealth)
	}
}
