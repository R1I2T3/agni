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

func InitializeMySQL(mySQLConfig db.MySQLConfig) {
	envConfig := GetEnvConfig()
	allModel := []interface{}{
		&db.Application{},
		&db.Notification{},
		&db.WebPushSubscription{},
	}
	if err := db.InitMySQL(envConfig.ENV_MODE, mySQLConfig, allModel...); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	mySQLHealth := db.MySQLHealthCheck()
	log.Printf("Database Health Status: %+v", mySQLHealth)

	if mySQLHealth["ping"] == true {
		log.Println("✅ Database is healthy")
	} else {
		log.Fatalf("❌ Database health check failed: %v", mySQLHealth)
	}
}
