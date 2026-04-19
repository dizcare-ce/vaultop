// Package quota enforces per-key operation limits within a rolling window.
package quota

import (
	"errors"
	"sync"
	"time"
)

// ErrQuotaExceeded is returned when a key has exceeded its allowed operations.
var ErrQuotaExceeded = errors.New("quota exceeded")

// Config holds quota settings.
type Config struct {
	MaxOps   int
	Window   time.Duration
}

type entry struct {
	count     int
	windowEnd time.Time
}

// Limiter tracks operation counts per key.
type Limiter struct {
	mu      sync.Mutex
	cfg     Config
	entries map[string]*entry
	now     func() time.Time
}

// New creates a Limiter with the given config.
func New(cfg Config) *Limiter {
	return &Limiter{
		cfg:     cfg,
		entries: make(map[string]*entry),
		now:     time.Now,
	}
}

// Allow checks and records an operation for key. Returns ErrQuotaExceeded if limit is reached.
func (l *Limiter) Allow(key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	e, ok := l.entries[key]
	if !ok || now.After(e.windowEnd) {
		l.entries[key] = &entry{count: 1, windowEnd: now.Add(l.cfg.Window)}
		return nil
	}
	if e.count >= l.cfg.MaxOps {
		return ErrQuotaExceeded
	}
	e.count++
	return nil
}

// Remaining returns how many operations are left for key in the current window.
func (l *Limiter) Remaining(key string) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	e, ok := l.entries[key]
	if !ok || now.After(e.windowEnd) {
		return l.cfg.MaxOps
	}
	return l.cfg.MaxOps - e.count
}

// Reset clears the quota state for a key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, key)
}
