package main // Main entry point for the worker service.

import (
	"fmt"
	"net/http"

	"api-gateway/config"
	"api-gateway/internal"
)

func main() {
	cfg := config.LoadConfig()
	router := &internal.Router{TaskServiceURL: cfg.TaskServiceURL}

	fmt.Println("API Gateway is running on port 9090")
	http.ListenAndServe(":9090", router)
}
