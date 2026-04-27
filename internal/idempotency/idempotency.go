// Package idempotency provides request deduplication for secret operations.
// It ensures that retried operations with the same key produce the same
// outcome without applying side-effects more than once.
package idempotency

import (
	"errors"
	"sync"
	"time"
)

// ErrDuplicate is returned when an operation with the given key has already
// been completed within the active window.
var ErrDuplicate = errors.New("idempotency: duplicate request")

// Result holds the outcome of a completed operation.
type Result struct {
	Value     string
	Err       error
	CompletedAt time.Time
}

// entry is an internal record stored per idempotency key.
type entry struct {
	result    Result
	expiresAt time.Time
}

// Store deduplicates operations within a configurable time window.
type Store struct {
	mu      sync.Mutex
	records map[string]entry
	window  time.Duration
	clock   func() time.Time
}

// New returns a Store with the given deduplication window.
func New(window time.Duration) *Store {
	return newWithClock(window, time.Now)
}

func newWithClock(window time.Duration, clock func() time.Time) *Store {
	return &Store{
		records: make(map[string]entry),
		window:  window,
		clock:   clock,
	}
}

// Check returns (result, true) if the key was already completed within the
// window, or (zero, false) if the key is new or has expired.
func (s *Store) Check(key string) (Result, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.clock()
	if e, ok := s.records[key]; ok && now.Before(e.expiresAt) {
		return e.result, true
	}
	return Result{}, false
}

// Record stores the result of a completed operation under the given key.
func (s *Store) Record(key string, r Result) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[key] = entry{
		result:    r,
		expiresAt: s.clock().Add(s.window),
	}
}

// Purge removes all expired entries.
func (s *Store) Purge() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.clock()
	for k, e := range s.records {
		if !now.Before(e.expiresAt) {
			delete(s.records, k)
		}
	}
}

// Len returns the number of active (non-expired) entries.
func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.clock()
	count := 0
	for _, e := range s.records {
		if now.Before(e.expiresAt) {
			count++
		}
	}
	return count
}
