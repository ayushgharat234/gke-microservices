package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"notification-service/internal/email"
	"notification-service/internal/slack"
)

// NotificationRequest defines the structure of the incoming request
// It now includes a 'Type' field to specify the notification channel.
type NotificationRequest struct {
	Type    string `json:"type"`
	TaskID  string `json:"task_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Webhook string `json:"webhook_url"` // Webhook URL for Slack
	EmailTo string `json:"email_to"`    // Recipient email address
}

// NotifyHandler handles notifications for task status updates,
// routing them to the correct service (Slack, email, etc.)
func NotifyHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req NotificationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Notification received for Task [%s]: Status [%s], Message [%s]\n", req.TaskID, req.Status, req.Message)

	// Build the notification message string
	notifMessage := "Task: " + req.TaskID + "\nStatus: " + req.Status + "\nMessage: " + req.Message

	switch req.Type {
	case "slack":
		if req.Webhook == "" {
			http.Error(w, "Slack webhook URL is required", http.StatusBadRequest)
			return
		}
		if err := slack.SendSlackMessage(req.Webhook, notifMessage); err != nil {
			log.Printf("Failed to send Slack message: %v", err)
			http.Error(w, "Failed to send Slack message", http.StatusInternalServerError)
			return
		}
	case "email":
		if req.EmailTo == "" {
			http.Error(w, "Recipient email is required", http.StatusBadRequest)
			return
		}
		if err := email.SendEmail(req.EmailTo, "Task Status Update", notifMessage); err != nil {
			log.Printf("Failed to send email: %v", err)
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}
	case "both":
		// Send to Slack first
		if err := slack.SendSlackMessage(req.Webhook, notifMessage); err != nil {
			log.Printf("Failed to send Slack message: %v", err)
		}
		// Send to Email
		if err := email.SendEmail(req.EmailTo, "Task Status Update", notifMessage); err != nil {
			log.Printf("Failed to send email: %v", err)
		}
	default:
		http.Error(w, "Invalid notification type. Use 'slack', 'email', or 'both'", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification processed successfully"))
}
