package internal // Package internal contains the core logic for the worker service.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

// Struct for the TaskServiceHandler
type Task struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"` // pending, in-progress, completed
}

// Worker for Polling the Queeue of Tasks
func StartWorker(client *redis.Client) {
	ctx := context.Background()

	for {
		result, err := client.LPop(ctx, "task_queue").Result()
		if err == redis.Nil {
			fmt.Println("No tasks in queue, waiting...")
			time.Sleep(5 * time.Second)
			continue
		} else if err != nil {
			fmt.Printf("Error popping task from queue: %s\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		var task Task
		if err := json.Unmarshal([]byte(result), &task); err != nil {
			fmt.Printf("Error unmarshalling task: %s\n", err)
			// In a production system, you might move this message to a dead-letter queue
			continue
		}

		fmt.Printf("Processing task [%s]: %s\n", task.ID, task.Title)
		time.Sleep(2 * time.Second) // Simulate task processing
		fmt.Printf("Completed task [%s]\n", task.ID)
		
		// Send notification to notification service
		sendNotificationToService(task.ID, "completed", "Task has been processed successfully")
	}
}

// sendNotificationToService sends a notification to the notification service
func sendNotificationToService(taskID, status, message string) {
	notif := map[string]string{
		"task_id": taskID,
		"status":  status,
		"message": message,
	}

	data, _ := json.Marshal(notif)
	_, err := http.Post("http://notification-service:8083/notify", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
	} else {
		log.Printf("Notification sent successfully for task [%s]", taskID)
	}
}

