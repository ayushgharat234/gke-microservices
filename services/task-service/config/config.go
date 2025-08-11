package config // Package config provides configuration settings for the task service.

import (
	"os"
)

// Struct for configuration settings
type Config struct {
	RedisAddr string
	RedisPass string
}

func LoadConfig() *Config {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPass := os.Getenv("REDIS_PASS")

	if redisAddr == "" {
		redisAddr = "localhost:6379" // Default Redis address
	}

	return &Config{
		RedisAddr: redisAddr,
		RedisPass: redisPass,
	}
}
