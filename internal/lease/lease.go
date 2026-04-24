// Package lease provides time-bounded secret leases with automatic expiry
// tracking. A lease associates a secret key with a duration, an issue time,
// and an optional renewal count so callers can enforce maximum lifetimes.
package lease

import (
	"errors"
	"sync"
	"time"
)

// ErrNotFound is returned when a lease does not exist for the given key.
var ErrNotFound = errors.New("lease: not found")

// ErrExpired is returned when the lease exists but has already expired.
var ErrExpired = errors.New("lease: expired")

// Lease describes a single time-bounded secret lease.
type Lease struct {
	Key       string
	IssuedAt  time.Time
	Duration  time.Duration
	Renewals  int
}

// ExpiresAt returns the absolute expiry time for the lease.
func (l Lease) ExpiresAt() time.Time {
	return l.IssuedAt.Add(l.Duration)
}

// IsExpired reports whether the lease has expired relative to now.
func (l Lease) IsExpired(now time.Time) bool {
	return now.After(l.ExpiresAt())
}

// Manager stores and manages leases in memory.
type Manager struct {
	mu     sync.Mutex
	leases map[string]Lease
	now    func() time.Time
}

// New returns a new Manager using the real clock.
func New() *Manager {
	return &Manager{
		leases: make(map[string]Lease),
		now:    time.Now,
	}
}

// Grant creates or replaces the lease for key with the given duration.
func (m *Manager) Grant(key string, d time.Duration) Lease {
	m.mu.Lock()
	defer m.mu.Unlock()
	l := Lease{Key: key, IssuedAt: m.now(), Duration: d}
	m.leases[key] = l
	return l
}

// Renew extends an existing lease by adding d to its remaining time.
// Returns ErrNotFound if the lease does not exist.
func (m *Manager) Renew(key string, d time.Duration) (Lease, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.leases[key]
	if !ok {
		return Lease{}, ErrNotFound
	}
	now := m.now()
	remaining := l.ExpiresAt().Sub(now)
	if remaining < 0 {
		remaining = 0
	}
	l.IssuedAt = now
	l.Duration = remaining + d
	l.Renewals++
	m.leases[key] = l
	return l, nil
}

// Get returns the lease for key, or an error if missing or expired.
func (m *Manager) Get(key string) (Lease, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.leases[key]
	if !ok {
		return Lease{}, ErrNotFound
	}
	if l.IsExpired(m.now()) {
		return l, ErrExpired
	}
	return l, nil
}

// Revoke removes the lease for key. It is a no-op if the key is unknown.
func (m *Manager) Revoke(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.leases, key)
}

// Active returns all leases that have not yet expired.
func (m *Manager) Active() []Lease {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.now()
	out := make([]Lease, 0, len(m.leases))
	for _, l := range m.leases {
		if !l.IsExpired(now) {
			out = append(out, l)
		}
	}
	return out
}

// Purge removes all expired leases and returns the count removed.
func (m *Manager) Purge() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.now()
	n := 0
	for k, l := range m.leases {
		if l.IsExpired(now) {
			delete(m.leases, k)
			n++
		}
	}
	return n
}
