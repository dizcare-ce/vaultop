// Package ratelimit provides a simple token-bucket rate limiter for
// controlling how frequently secrets may be rotated or accessed.
package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// ErrRateLimited is returned when the rate limit for a key is exceeded.
var ErrRateLimited = fmt.Errorf("rate limit exceeded")

// Limiter tracks per-key request counts within a rolling window.
type Limiter struct {
	mu       sync.Mutex
	window   time.Duration
	max      int
	buckets  map[string][]time.Time
	now      func() time.Time
}

// New creates a Limiter that allows at most max calls per window per key.
func New(max int, window time.Duration) *Limiter {
	return &Limiter{
		max:     max,
		window:  window,
		buckets: make(map[string][]time.Time),
		now:     time.Now,
	}
}

// Allow returns nil if the call is within the rate limit, or ErrRateLimited.
func (l *Limiter) Allow(key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	cutoff := now.Add(-l.window)

	times := l.buckets[key]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= l.max {
		l.buckets[key] = filtered
		return ErrRateLimited
	}

	l.buckets[key] = append(filtered, now)
	return nil
}

// Reset clears the call history for a key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, key)
}

// Remaining returns how many calls are still allowed for key within the window.
func (l *Limiter) Remaining(key string) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	cutoff := now.Add(-l.window)
	count := 0
	for _, t := range l.buckets[key] {
		if t.After(cutoff) {
			count++
		}
	}
	rem := l.max - count
	if rem < 0 {
		return 0
	}
	return rem
}
