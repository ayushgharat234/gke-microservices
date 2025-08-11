package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

// Struct for the Notification
type Notification struct {
	TaskID  string `json:"task_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// NotifyHandler handles notifications for task status updates
func NotifyHandler(w http.ResponseWriter, r *http.Request) {
	var notif Notification
	err := json.NewDecoder(r.Body).Decode(&notif)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Notification received for Task [%s]: Status [%s], Message [%s]\n", notif.TaskID, notif.Status, notif.Message)
	w.WriteHeader(http.StatusOK)
}
