package config

import (
	"os"
	"testing"
	"time"
)

func TestGetEnvAsDuration(t *testing.T) {
	// Clean up environment at the end
	defer os.Unsetenv("TEST_DURATION")

	// 1. Default value parsing
	dur := GetEnvAsDuration("TEST_DURATION", 10*time.Second)
	if dur != 10*time.Second {
		t.Errorf("expected 10s, got %v", dur)
	}

	// 2. Correct string format parsing
	os.Setenv("TEST_DURATION", "5s")
	dur = GetEnvAsDuration("TEST_DURATION", 10*time.Second)
	if dur != 5*time.Second {
		t.Errorf("expected 5s, got %v", dur)
	}

	// 3. Fallback on invalid format
	os.Setenv("TEST_DURATION", "invalid")
	dur = GetEnvAsDuration("TEST_DURATION", 10*time.Second)
	if dur != 10*time.Second {
		t.Errorf("expected 10s on invalid format, got %v", dur)
	}
}

func TestENV_MODE_Defaults(t *testing.T) {
	defer os.Unsetenv("ENV_MODE")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("REDIS_HOST")

	// 1. Default should be "prod"
	os.Unsetenv("ENV_MODE")
	cfg := GetEnvConfig()
	if cfg.ENV_MODE != "prod" {
		t.Errorf("expected default ENV_MODE to be 'prod', got '%s'", cfg.ENV_MODE)
	}

	// 2. MySQL DB_HOST in prod mode should default to "agni-mysql"
	mysqlCfg := GetMySQLDBConfig()
	if mysqlCfg.DB_HOST != "agni-mysql" {
		t.Errorf("expected prod DB_HOST to be 'agni-mysql', got '%s'", mysqlCfg.DB_HOST)
	}

	// 3. Redis host in prod mode should default to "agni-redis"
	redisCfg := GetRedisEnvConfig()
	if redisCfg.Host != "agni-redis" {
		t.Errorf("expected prod Redis host to be 'agni-redis', got '%s'", redisCfg.Host)
	}

	// 4. In "local" mode, MySQL DB_HOST should default to "localhost"
	os.Setenv("ENV_MODE", "local")
	mysqlCfgLocal := GetMySQLDBConfig()
	if mysqlCfgLocal.DB_HOST != "localhost" {
		t.Errorf("expected local DB_HOST to be 'localhost', got '%s'", mysqlCfgLocal.DB_HOST)
	}

	// 5. In "local" mode, Redis host should default to "localhost"
	redisCfgLocal := GetRedisEnvConfig()
	if redisCfgLocal.Host != "localhost" {
		t.Errorf("expected local Redis host to be 'localhost', got '%s'", redisCfgLocal.Host)
	}

	// 6. Explicit env overrides should still be respected
	os.Setenv("DB_HOST", "my-custom-mysql")
	os.Setenv("REDIS_HOST", "my-custom-redis")
	
	mysqlCfgOverride := GetMySQLDBConfig()
	if mysqlCfgOverride.DB_HOST != "my-custom-mysql" {
		t.Errorf("expected overridden DB_HOST to be 'my-custom-mysql', got '%s'", mysqlCfgOverride.DB_HOST)
	}

	redisCfgOverride := GetRedisEnvConfig()
	if redisCfgOverride.Host != "my-custom-redis" {
		t.Errorf("expected overridden Redis host to be 'my-custom-redis', got '%s'", redisCfgOverride.Host)
	}
}
