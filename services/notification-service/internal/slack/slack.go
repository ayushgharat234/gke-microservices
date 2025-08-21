package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SlackWebhookPayload represents the payload structure for a Slack webhook.
type SlackWebhookPayload struct {
	Text string `json:"text"`
}

// SendSlackMessage sends a message to a Slack channel via a webhook URL.
func SendSlackMessage(webhookURL, message string) error {
	// Create the JSON payload.
	payload := SlackWebhookPayload{
		Text: message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack payload: %w", err)
	}

	// Send the POST request to the Slack webhook.
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send Slack message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Slack returned non-OK status: %s", resp.Status)
	}

	return nil
}
