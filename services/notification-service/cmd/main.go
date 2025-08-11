package main

import (
	"log"
	"net/http"
	"notification-service/internal"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/health", http.StatusSeeOther)
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Notification Service is healthy"))
	})
	http.HandleFunc("/notify", internal.NotifyHandler)

	log.Printf("Notification Service is running on :8083")
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
