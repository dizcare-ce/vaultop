// Package history tracks rotation events for secrets,
// persisting a lightweight JSON log so schedules and audits
// can query when a secret was last rotated.
package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry records a single completed rotation.
type Entry struct {
	SecretKey  string    `json:"secret_key"`
	RotatedAt  time.Time `json:"rotated_at"`
	Provider   string    `json:"provider"`
	Success    bool      `json:"success"`
}

// Store is a simple file-backed rotation history store.
type Store struct {
	mu      sync.RWMutex
	path    string
	entries []Entry
}

// Load opens (or creates) the history file at path.
func Load(path string) (*Store, error) {
	s := &Store{path: path}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &s.entries); err != nil {
		return nil, err
	}
	return s, nil
}

// Record appends an entry and persists the store.
func (s *Store) Record(e Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, e)
	return s.flush()
}

// LastRotated returns the most recent successful rotation time for key.
// Returns zero time if no record exists.
func (s *Store) LastRotated(key string) time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var latest time.Time
	for _, e := range s.entries {
		if e.SecretKey == key && e.Success && e.RotatedAt.After(latest) {
			latest = e.RotatedAt
		}
	}
	return latest
}

// All returns a copy of all entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	return out
}

func (s *Store) flush() error {
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}
