// Package version provides secret version tracking and retrieval.
// It records each write to a secret as a numbered version and allows
// callers to list or retrieve historical values.
package version

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrNotFound is returned when a key or version does not exist.
var ErrNotFound = errors.New("version: not found")

// Entry represents a single versioned value of a secret.
type Entry struct {
	Version   int
	Value     string
	CreatedAt time.Time
}

// Store holds versioned history for secrets in memory.
type Store struct {
	mu      sync.RWMutex
	records map[string][]Entry
}

// New creates an empty version Store.
func New() *Store {
	return &Store{records: make(map[string][]Entry)}
}

// Record appends a new version for key with the given value.
// Versions are numbered starting at 1.
func (s *Store) Record(key, value string) Entry {
	s.mu.Lock()
	defer s.mu.Unlock()

	versions := s.records[key]
	e := Entry{
		Version:   len(versions) + 1,
		Value:     value,
		CreatedAt: time.Now().UTC(),
	}
	s.records[key] = append(versions, e)
	return e
}

// Get returns the Entry for the given key and version number.
func (s *Store) Get(key string, version int) (Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	versions, ok := s.records[key]
	if !ok || version < 1 || version > len(versions) {
		return Entry{}, fmt.Errorf("%w: key=%s version=%d", ErrNotFound, key, version)
	}
	return versions[version-1], nil
}

// Latest returns the most recent Entry for key.
func (s *Store) Latest(key string) (Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	versions, ok := s.records[key]
	if !ok || len(versions) == 0 {
		return Entry{}, fmt.Errorf("%w: key=%s", ErrNotFound, key)
	}
	return versions[len(versions)-1], nil
}

// List returns all recorded entries for key in ascending version order.
func (s *Store) List(key string) []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	versions := s.records[key]
	result := make([]Entry, len(versions))
	copy(result, versions)
	return result
}

// Count returns the total number of versions stored for key.
func (s *Store) Count(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.records[key])
}
