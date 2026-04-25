// Package deadline tracks per-key operation deadlines and reports
// violations when an operation exceeds its allowed duration.
package deadline

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrExceeded is returned when an operation has exceeded its deadline.
var ErrExceeded = errors.New("deadline exceeded")

// Entry holds the deadline state for a single key.
type Entry struct {
	Key       string
	Deadline  time.Time
	StartedAt time.Time
}

// Exceeded reports whether the entry's deadline has passed relative to now.
func (e Entry) Exceeded(now time.Time) bool {
	return now.After(e.Deadline)
}

// Remaining returns the duration left before the deadline.
func (e Entry) Remaining(now time.Time) time.Duration {
	d := e.Deadline.Sub(now)
	if d < 0 {
		return 0
	}
	return d
}

// Tracker manages deadlines for concurrent operations.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]Entry
	clock   func() time.Time
}

// New returns a new Tracker using the real wall clock.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]Entry),
		clock:   time.Now,
	}
}

// Start registers a deadline for key, expiring after ttl from now.
// Calling Start for an already-tracked key overwrites the previous entry.
func (t *Tracker) Start(key string, ttl time.Duration) {
	now := t.clock()
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[key] = Entry{
		Key:       key,
		Deadline:  now.Add(ttl),
		StartedAt: now,
	}
}

// Finish removes the deadline for key. It is a no-op if the key is unknown.
func (t *Tracker) Finish(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, key)
}

// Check returns ErrExceeded if the deadline for key has passed, nil if it is
// still within bounds, or an error if the key is not being tracked.
func (t *Tracker) Check(key string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[key]
	if !ok {
		return fmt.Errorf("deadline: key %q not tracked", key)
	}
	if e.Exceeded(t.clock()) {
		return fmt.Errorf("%w: key %q exceeded at %s", ErrExceeded, key, e.Deadline.Format(time.RFC3339))
	}
	return nil
}

// Active returns all entries that have not yet exceeded their deadline.
func (t *Tracker) Active() []Entry {
	now := t.clock()
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		if !e.Exceeded(now) {
			out = append(out, e)
		}
	}
	return out
}

// Violations returns all entries whose deadline has already passed.
func (t *Tracker) Violations() []Entry {
	now := t.clock()
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, 0)
	for _, e := range t.entries {
		if e.Exceeded(now) {
			out = append(out, e)
		}
	}
	return out
}
