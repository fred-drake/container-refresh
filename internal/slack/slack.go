package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type SlackMessage struct {
	Text string `json:"text"`
}

// SendMessage sends a message to Slack using the webhook URL from the environment
func SendMessage(message string) error {
	webhookURLFile := os.Getenv("SLACK_WEBHOOK_URL_FILE")
	if webhookURLFile == "" {
		return nil // Skip if webhook URL file is not configured
	}

	webhookURLBytes, err := os.ReadFile(webhookURLFile)
	if err != nil {
		return fmt.Errorf("failed to read Slack webhook URL file: %v", err)
	}
	webhookURL := strings.TrimSpace(string(webhookURLBytes))

	payload := SlackMessage{Text: message}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack message: %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send Slack message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send Slack message: received status code %d", resp.StatusCode)
	}

	return nil
}
