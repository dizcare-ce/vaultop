// Package observe provides a lightweight observability hook that aggregates
// per-key operation counters, latency samples, and last-seen timestamps so
// that operators can inspect runtime behaviour without an external metrics
// backend.
package observe

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Entry holds aggregated observations for a single key.
type Entry struct {
	Key       string
	Ops       int64
	Errors    int64
	TotalNs   int64 // cumulative latency in nanoseconds
	LastSeen  time.Time
}

// AvgLatency returns the mean operation latency, or zero if no ops recorded.
func (e Entry) AvgLatency() time.Duration {
	if e.Ops == 0 {
		return 0
	}
	return time.Duration(e.TotalNs / e.Ops)
}

// Observer records observations for named keys.
type Observer struct {
	mu      sync.Mutex
	entries map[string]*Entry
	clock   func() time.Time
}

// New returns a new Observer.
func New() *Observer {
	return &Observer{
		entries: make(map[string]*Entry),
		clock:   time.Now,
	}
}

// Record registers a completed operation for key with the given latency and
// whether it resulted in an error.
func (o *Observer) Record(key string, latency time.Duration, err error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	e, ok := o.entries[key]
	if !ok {
		e = &Entry{Key: key}
		o.entries[key] = e
	}
	e.Ops++
	e.TotalNs += latency.Nanoseconds()
	e.LastSeen = o.clock()
	if err != nil {
		e.Errors++
	}
}

// Get returns the Entry for key and whether it was found.
func (o *Observer) Get(key string) (Entry, bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	e, ok := o.entries[key]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a snapshot of all recorded entries.
func (o *Observer) All() []Entry {
	o.mu.Lock()
	defer o.mu.Unlock()
	out := make([]Entry, 0, len(o.entries))
	for _, e := range o.entries {
		out = append(out, *e)
	}
	return out
}

// WriteSummary writes a human-readable summary table to w.
func (o *Observer) WriteSummary(w io.Writer) {
	for _, e := range o.All() {
		fmt.Fprintf(w, "key=%-30s ops=%-6d errors=%-4d avg_latency=%s last_seen=%s\n",
			e.Key, e.Ops, e.Errors, e.AvgLatency().Round(time.Microsecond),
			e.LastSeen.Format(time.RFC3339))
	}
}
