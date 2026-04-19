package webhook

import (
	"context"
	"log"
)

// Dispatcher fans out a single event to multiple Notifiers.
type Dispatcher struct {
	targets []Notifier
	logger  *log.Logger
}

// NewDispatcher returns a Dispatcher that delivers to all provided targets.
func NewDispatcher(logger *log.Logger, targets ...Notifier) *Dispatcher {
	return &Dispatcher{targets: targets, logger: logger}
}

// Send delivers the event to every target. Errors are logged but do not
// prevent delivery to remaining targets. Returns the last error encountered.
func (d *Dispatcher) Send(ctx context.Context, e Event) error {
	var last error
	for _, t := range d.targets {
		if err := t.Send(ctx, e); err != nil {
			if d.logger != nil {
				d.logger.Printf("webhook dispatcher: %v", err)
			}
			last = err
		}
	}
	return last
}
