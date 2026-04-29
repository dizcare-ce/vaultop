// Package window provides a sliding-window counter for tracking event
// frequency over a rolling time period.
package window

import (
	"sync"
	"time"
)

// entry records a single event timestamp.
type entry struct {
	at time.Time
}

// Counter is a thread-safe sliding-window event counter keyed by name.
type Counter struct {
	mu     sync.Mutex
	window time.Duration
	buckets map[string][]entry
	clock  func() time.Time
}

// New returns a Counter that tracks events within the given rolling window.
func New(window time.Duration) *Counter {
	return newWithClock(window, time.Now)
}

func newWithClock(window time.Duration, clock func() time.Time) *Counter {
	return &Counter{
		window:  window,
		buckets: make(map[string][]entry),
		clock:   clock,
	}
}

// Record adds one event for the given key at the current time.
func (c *Counter) Record(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clock()
	c.buckets[key] = append(c.evict(key, now), entry{at: now})
}

// Count returns the number of events for key within the current window.
func (c *Counter) Count(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clock()
	c.buckets[key] = c.evict(key, now)
	return len(c.buckets[key])
}

// Reset clears all recorded events for the given key.
func (c *Counter) Reset(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.buckets, key)
}

// evict removes entries that have fallen outside the window. Caller must hold mu.
func (c *Counter) evict(key string, now time.Time) []entry {
	cutoff := now.Add(-c.window)
	old := c.buckets[key]
	var fresh []entry
	for _, e := range old {
		if !e.at.Before(cutoff) {
			fresh = append(fresh, e)
		}
	}
	return fresh
}
