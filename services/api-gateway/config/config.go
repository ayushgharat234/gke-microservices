package config // Package config provides configuration settings for the task service.

import "os"

// Config holds the configuration settings for the task service.
type Config struct {
	TaskServiceURL string
}

// LoadConfig loads the configuration from environment variables or sets defaults.
func LoadConfig() *Config {
	url := os.Getenv("TASK_SERVICE_URL")
	if url == "" {
		url = "http://localhost:9090" // Fallback
	}

	return &Config{TaskServiceURL: url}
}
