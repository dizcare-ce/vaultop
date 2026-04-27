// Package event provides a simple in-process event bus for broadcasting
// secret lifecycle events to registered subscribers.
package event

import "sync"

// Kind identifies the type of event.
type Kind string

const (
	KindRotated  Kind = "rotated"
	KindDeleted  Kind = "deleted"
	KindImported Kind = "imported"
	KindExpired  Kind = "expired"
)

// Event carries information about a secret lifecycle occurrence.
type Event struct {
	Kind    Kind
	Key     string
	Message string
}

// Handler is a function that receives an event.
type Handler func(Event)

// Bus dispatches events to registered handlers.
type Bus struct {
	mu       sync.RWMutex
	handlers map[Kind][]Handler
}

// New returns an initialised Bus.
func New() *Bus {
	return &Bus{handlers: make(map[Kind][]Handler)}
}

// Subscribe registers h to be called whenever an event of the given kind is
// published. Subscribing with KindAll (empty string) is not supported; use
// explicit kinds.
func (b *Bus) Subscribe(kind Kind, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[kind] = append(b.handlers[kind], h)
}

// Publish sends e to every handler registered for e.Kind. Handlers are called
// synchronously in subscription order.
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	handlers := make([]Handler, len(b.handlers[e.Kind]))
	copy(handlers, b.handlers[e.Kind])
	b.mu.RUnlock()

	for _, h := range handlers {
		h(e)
	}
}

// SubscriberCount returns the number of handlers registered for the given kind.
func (b *Bus) SubscriberCount(kind Kind) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.handlers[kind])
}

// Reset removes all subscriptions.
func (b *Bus) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = make(map[Kind][]Handler)
}
