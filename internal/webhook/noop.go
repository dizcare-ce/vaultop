package webhook

import "context"

// Notifier is the interface for sending webhook events.
type Notifier interface {
	Send(ctx context.Context, e Event) error
}

// Noop is a no-op Notifier that discards all events.
type Noop struct{}

// Send does nothing and returns nil.
func (Noop) Send(_ context.Context, _ Event) error { return nil }
