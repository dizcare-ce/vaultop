// Package notification provides interfaces and implementations for
// notifying operators when secrets are rotated or rotation fails.
package notification

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Event represents a rotation notification event.
type Event struct {
	SecretKey string
	Provider  string
	Success   bool
	Error     error
	RotatedAt time.Time
}

// Notifier sends rotation events to an output destination.
type Notifier interface {
	Notify(e Event) error
}

// WriterNotifier writes human-readable event lines to an io.Writer.
type WriterNotifier struct {
	w io.Writer
}

// NewWriterNotifier returns a Notifier that writes to w.
func NewWriterNotifier(w io.Writer) *WriterNotifier {
	return &WriterNotifier{w: w}
}

// NewStdoutNotifier returns a Notifier that writes to stdout.
func NewStdoutNotifier() *WriterNotifier {
	return NewWriterNotifier(os.Stdout)
}

// Notify formats and writes the event to the underlying writer.
func (n *WriterNotifier) Notify(e Event) error {
	status := "OK"
	detail := ""
	if !e.Success {
		status = "FAILED"
		if e.Error != nil {
			detail = fmt.Sprintf(" error=%q", e.Error.Error())
		}
	}
	_, err := fmt.Fprintf(
		n.w,
		"%s provider=%s key=%s status=%s%s\n",
		e.RotatedAt.UTC().Format(time.RFC3339),
		e.Provider,
		e.SecretKey,
		status,
		detail,
	)
	return err
}

// Noop is a Notifier that discards all events.
type Noop struct{}

// Notify does nothing and returns nil.
func (Noop) Notify(_ Event) error { return nil }
