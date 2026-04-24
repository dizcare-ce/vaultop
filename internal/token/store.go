package token

import (
	"errors"
	"sync"
	"time"
)

// ErrRevoked is returned when a token has been explicitly revoked.
var ErrRevoked = errors.New("token: revoked")

// Store tracks issued tokens and supports revocation.
type Store struct {
	mu      sync.RWMutex
	revoked map[string]struct{}
	issued  map[string]time.Time // value -> expiry
}

// NewStore returns an initialised Store.
func NewStore() *Store {
	return &Store{
		revoked: make(map[string]struct{}),
		issued:  make(map[string]time.Time),
	}
}

// Track records an issued token so it can later be revoked or purged.
func (s *Store) Track(t Token) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.issued[t.Value] = t.ExpiresAt
}

// Revoke marks a token as invalid regardless of its expiry.
func (s *Store) Revoke(value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.revoked[value] = struct{}{}
}

// IsRevoked reports whether the token has been revoked.
func (s *Store) IsRevoked(value string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.revoked[value]
	return ok
}

// Purge removes expired tokens from the store to reclaim memory.
func (s *Store) Purge(now time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	removed := 0
	for v, exp := range s.issued {
		if now.After(exp) {
			delete(s.issued, v)
			delete(s.revoked, v)
			removed++
		}
	}
	return removed
}

// Count returns the number of currently tracked (non-purged) tokens.
func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.issued)
}
