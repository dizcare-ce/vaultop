// Package audit provides a simple audit log for secret rotation events.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// EventType describes the kind of audit event.
type EventType string

const (
	EventRotated EventType = "rotated"
	EventDryRun  EventType = "dry_run"
	EventFailed  EventType = "failed"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     EventType `json:"event"`
	Provider  string    `json:"provider"`
	SecretKey string    `json:"secret_key"`
	Message   string    `json:"message,omitempty"`
}

// Logger writes audit entries to an io.Writer as newline-delimited JSON.
type Logger struct {
	w io.Writer
}

// NewLogger creates a Logger that writes to w.
func NewLogger(w io.Writer) *Logger {
	return &Logger{w: w}
}

// NewFileLogger opens (or creates) a file at path for append-only audit logging.
func NewFileLogger(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return NewLogger(f), nil
}

// Log writes an audit entry.
func (l *Logger) Log(event EventType, provider, secretKey, message string) error {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Provider:  provider,
		SecretKey: secretKey,
		Message:   message,
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}
