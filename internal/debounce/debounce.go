// Package debounce provides a mechanism to suppress repeated operations
// within a configurable quiet period. This is useful when multiple secret
// change events arrive in quick succession and only the final state matters.
package debounce

import (
	"sync"
	"time"
)

// Func is a function that can be debounced.
type Func func(key string)

// Debouncer delays execution of a function until a quiet period has elapsed
// since the last call for a given key.
type Debouncer struct {
	mu      sync.Mutex
	delay   time.Duration
	timers  map[string]*time.Timer
	clock   func() time.Time
	sleep   func(d time.Duration)
}

// New creates a Debouncer with the given quiet-period delay.
func New(delay time.Duration) *Debouncer {
	return &Debouncer{
		delay:  delay,
		timers: make(map[string]*time.Timer),
		clock:  time.Now,
		sleep:  time.Sleep,
	}
}

// Trigger schedules fn to be called after the debounce delay for key.
// If Trigger is called again for the same key before the delay elapses,
// the previous scheduled call is cancelled and the timer resets.
func (d *Debouncer) Trigger(key string, fn Func) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
	}

	d.timers[key] = time.AfterFunc(d.delay, func() {
		d.mu.Lock()
		delete(d.timers, key)
		d.mu.Unlock()
		fn(key)
	})
}

// Cancel stops any pending debounced call for key.
// It is a no-op if no call is pending.
func (d *Debouncer) Cancel(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
		delete(d.timers, key)
	}
}

// Pending returns the set of keys that have a pending debounced call.
func (d *Debouncer) Pending() []string {
	d.mu.Lock()
	defer d.mu.Unlock()

	keys := make([]string, 0, len(d.timers))
	for k := range d.timers {
		keys = append(keys, k)
	}
	return keys
}
