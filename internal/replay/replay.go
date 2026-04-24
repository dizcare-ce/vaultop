// Package replay provides protection against duplicate secret operation
// requests by tracking recently seen operation IDs within a sliding window.
package replay

import (
	"errors"
	"sync"
	"time"
)

// ErrReplay is returned when an operation ID has already been processed.
var ErrReplay = errors.New("replay: duplicate operation ID")

// entry holds the expiry time for a seen operation ID.
type entry struct {
	expiresAt time.Time
}

// Detector tracks operation IDs to detect replayed requests.
type Detector struct {
	mu      sync.Mutex
	seen    map[string]entry
	window  time.Duration
	clock   func() time.Time
}

// New creates a Detector with the given replay window duration.
// Operation IDs are remembered for the duration of the window.
func New(window time.Duration) *Detector {
	return &Detector{
		seen:   make(map[string]entry),
		window: window,
		clock:  time.Now,
	}
}

// Check returns ErrReplay if the given operation ID has been seen within the
// current window. Otherwise it records the ID and returns nil.
func (d *Detector) Check(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.clock()
	d.evict(now)

	if _, exists := d.seen[id]; exists {
		return ErrReplay
	}

	d.seen[id] = entry{expiresAt: now.Add(d.window)}
	return nil
}

// Seen reports whether the given ID is currently tracked (not yet expired).
func (d *Detector) Seen(id string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.clock()
	e, ok := d.seen[id]
	return ok && e.expiresAt.After(now)
}

// evict removes all expired entries. Must be called with d.mu held.
func (d *Detector) evict(now time.Time) {
	for id, e := range d.seen {
		if !e.expiresAt.After(now) {
			delete(d.seen, id)
		}
	}
}

// Size returns the number of currently tracked (non-expired) IDs.
func (d *Detector) Size() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.evict(d.clock())
	return len(d.seen)
}
