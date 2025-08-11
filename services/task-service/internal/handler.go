package internal // Package internal provides the HTTP handlers for the task service.

// Importing the necessary libraries and packages
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Struct for the TaskServiceHandler
type Task struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"` // pending, in-progress, completed
}

// Redis client instance with the context package
var ctx = context.Background()
var redisClient *redis.Client

// Initialize Redis
func InitRedis(addr, pass string) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		panic(fmt.Sprintf("Redis connection failed: %s", err))
	}

	fmt.Println("Connected to Redis")
}

// HealthHandler checks the health of the service
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task Service is healthy"))
}

// ReadinessHandler checks if the service is ready to accept requests
func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	if err := redisClient.Ping(ctx).Err(); err != nil {
		http.Error(w, "Redis is not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

// CreateTaskHandler creates a new task
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	task.ID = uuid.New().String() // Generate a unique ID for the task
	task.Status = "pending"       // Default Status

	// Store in Redis
	taskJSON, _ := json.Marshal(task)
	if err := redisClient.RPush(ctx, "task_queue", taskJSON).Err(); err != nil {
		http.Error(w, "Failed to store task", http.StatusInternalServerError)
		return
	}

	fmt.Printf("New Task Created [%s] at %s: %v+\n", task.ID, time.Now().Format(time.RFC3339), task)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
