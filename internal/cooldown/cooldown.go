// Package cooldown enforces a minimum wait period between successive
// operations on the same key, preventing rapid repeated actions.
package cooldown

import (
	"errors"
	"sync"
	"time"
)

// ErrCooldown is returned when an operation is attempted before the
// cooldown period for that key has elapsed.
var ErrCooldown = errors.New("cooldown: operation not permitted, key is cooling down")

// Clock allows time to be injected for testing.
type Clock func() time.Time

// Manager tracks per-key cooldown state.
type Manager struct {
	mu       sync.Mutex
	period   time.Duration
	last     map[string]time.Time
	clock    Clock
}

// New creates a Manager that enforces the given cooldown period.
// All keys share the same period.
func New(period time.Duration) *Manager {
	return &Manager{
		period: period,
		last:   make(map[string]time.Time),
		clock:  time.Now,
	}
}

// newWithClock creates a Manager with an injectable clock (for tests).
func newWithClock(period time.Duration, clock Clock) *Manager {
	m := New(period)
	m.clock = clock
	return m
}

// Allow checks whether the key is permitted to proceed. If the key is
// still within its cooldown window, ErrCooldown is returned. Otherwise
// the key's last-seen timestamp is updated and nil is returned.
func (m *Manager) Allow(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := m.clock()
	if last, ok := m.last[key]; ok {
		if now.Sub(last) < m.period {
			return ErrCooldown
		}
	}
	m.last[key] = now
	return nil
}

// Reset clears the cooldown state for a key, allowing it to proceed
// immediately on the next call to Allow.
func (m *Manager) Reset(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.last, key)
}

// Remaining returns how much cooldown time is left for the given key.
// Returns 0 if the key is not cooling down.
func (m *Manager) Remaining(key string) time.Duration {
	m.mu.Lock()
	defer m.mu.Unlock()

	last, ok := m.last[key]
	if !ok {
		return 0
	}
	elapsed := m.clock().Sub(last)
	if elapsed >= m.period {
		return 0
	}
	return m.period - elapsed
}
