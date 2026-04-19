package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookNotifier posts JSON payloads to an HTTP endpoint.
type WebhookNotifier struct {
	url    string
	client *http.Client
}

// NewWebhookNotifier returns a WebhookNotifier targeting url.
func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{
		url:    url,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

type webhookPayload struct {
	Level string   `json:"level"`
	Title string   `json:"title"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags,omitempty"`
}

// Send marshals msg to JSON and POSTs it to the configured URL.
func (w *WebhookNotifier) Send(msg Message) error {
	payload := webhookPayload{
		Level: string(msg.Level),
		Title: msg.Title,
		Body:  msg.Body,
		Tags:  msg.Tags,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook marshal: %w", err)
	}
	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("webhook post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// Name returns the backend identifier.
func (w *WebhookNotifier) Name() string { return "webhook" }
