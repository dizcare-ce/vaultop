// Package metrics tracks rotation and secret operation counters.
package metrics

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Counter holds cumulative counts for a named operation.
type Counter struct {
	Name      string
	Successes int
	Failures  int
	LastRun   time.Time
}

// Tracker accumulates operation metrics in memory.
type Tracker struct {
	mu       sync.Mutex
	counters map[string]*Counter
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{counters: make(map[string]*Counter)}
}

func (t *Tracker) get(name string) *Counter {
	if c, ok := t.counters[name]; ok {
		return c
	}
	c := &Counter{Name: name}
	t.counters[name] = c
	return c
}

// RecordSuccess increments the success counter for name.
func (t *Tracker) RecordSuccess(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	c := t.get(name)
	c.Successes++
	c.LastRun = time.Now()
}

// RecordFailure increments the failure counter for name.
func (t *Tracker) RecordFailure(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	c := t.get(name)
	c.Failures++
	c.LastRun = time.Now()
}

// Get returns a copy of the Counter for name, and whether it exists.
func (t *Tracker) Get(name string) (Counter, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	c, ok := t.counters[name]
	if !ok {
		return Counter{}, false
	}
	return *c, true
}

// All returns a snapshot of every counter.
func (t *Tracker) All() []Counter {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Counter, 0, len(t.counters))
	for _, c := range t.counters {
		out = append(out, *c)
	}
	return out
}

// WriteSummary writes a human-readable summary to w.
func (t *Tracker) WriteSummary(w io.Writer) {
	for _, c := range t.All() {
		fmt.Fprintf(w, "%-40s ok=%-6d fail=%-6d last=%s\n",
			c.Name, c.Successes, c.Failures, c.LastRun.Format(time.RFC3339))
	}
}
