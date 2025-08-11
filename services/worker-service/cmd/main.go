package main // Main entry point for the worker service.

import (
	"fmt"

	"worker-service/config"
	"worker-service/internal"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.LoadConfig()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr, // Redis server address
		Password: cfg.RedisPass, // Redis password
		DB:       0,             // Default DB
	})

	fmt.Println("Worker Service Started")
	internal.StartWorker(client)
}
