package config // Package config provides configuration settings for the work service.

import "os"

type Config struct {
	RedisAddr string
	RedisPass string
}

func LoadConfig() *Config {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	return &Config{
		RedisAddr: addr,
		RedisPass: os.Getenv("REDIS_PASS"),
	}
}
