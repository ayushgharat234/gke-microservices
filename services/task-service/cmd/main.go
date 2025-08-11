package main

import (
	"fmt"
	"log"
	"net/http"

	"task-service/config"
	"task-service/internal"
)

func main() {
	cfg := config.LoadConfig()
	internal.InitRedis(cfg.RedisAddr, cfg.RedisPass)

	// Routing the HTTP handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/health", http.StatusSeeOther)
	})
	http.HandleFunc("/health", internal.HealthHandler)
	http.HandleFunc("/create-task", internal.CreateTaskHandler)
	http.HandleFunc("/readiness", internal.ReadinessHandler)

	fmt.Println("Task Service is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
