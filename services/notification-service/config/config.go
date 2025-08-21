package config

import (
	"os"
)

// Config holds the application configuration.
type Config struct {
	SlackWebhookURL string
	EmailUser       string
	EmailPassword   string
	EmailSMTPHost   string
	EmailSMTPPort   string
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() *Config {
	return &Config{
		SlackWebhookURL: os.Getenv("SLACK_WEBHOOK_URL"),
		EmailUser:       os.Getenv("EMAIL_USER"),
		EmailPassword:   os.Getenv("EMAIL_PASSWORD"),
		EmailSMTPHost:   os.Getenv("EMAIL_SMTP_HOST"),
		EmailSMTPPort:   os.Getenv("EMAIL_SMTP_PORT"),
	}
}
