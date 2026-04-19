// Package webhook delivers event notifications to HTTP endpoints.
package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Event represents a webhook payload sent remote endpoint.
type Event struct {
	Kind      string            `json:"kind"`
	Key       string            `json:"key"`
	Timestamp time.Time         `json:"timestamp"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// Sender delivers webhook events to a configured URL.
type Sender struct {
	URL     string
	Timeout time.Duration
	client  *http.Client
}

// New returns a Sender targeting the given URL.
func New(url string, timeout time.Duration) *Sender {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &Sender{
		URL:     url,
		Timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

// Send marshals the event and POSTs it to the configured URL.
func (s *Sender) Send(ctx context.Context, e Event) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	buf, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("webhook: marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.URL, bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("webhook: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}
