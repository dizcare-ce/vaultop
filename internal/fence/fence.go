// Package fence provides a fencing token mechanism to prevent stale
// writers from overwriting data when distributed locks are involved.
//
// A fencing token is a monotonically increasing integer issued each time
// a lock is acquired. Writers must present the token when performing an
// operation; if the token is older than the last accepted one the write
// is rejected.
package fence

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrStaleFence is returned when the presented token is outdated.
var ErrStaleFence = errors.New("fence: stale token")

// Token represents a fencing token issued at a point in time.
type Token struct {
	Seq       uint64
	IssuedAt  time.Time
	HolderID  string
}

func (t Token) String() string {
	return fmt.Sprintf("fence(%d, holder=%s)", t.Seq, t.HolderID)
}

// Manager issues and validates fencing tokens per named resource.
type Manager struct {
	mu      sync.Mutex
	current map[string]Token
	clock   func() time.Time
}

// New returns a new Manager.
func New() *Manager {
	return &Manager{
		current: make(map[string]Token),
		clock:   time.Now,
	}
}

// Issue issues a new fencing token for the named resource and holder.
// Each call increments the sequence number for that resource.
func (m *Manager) Issue(resource, holderID string) Token {
	m.mu.Lock()
	defer m.mu.Unlock()

	prev := m.current[resource]
	tok := Token{
		Seq:      prev.Seq + 1,
		IssuedAt: m.clock(),
		HolderID: holderID,
	}
	m.current[resource] = tok
	return tok
}

// Check validates that tok is still the current token for resource.
// Returns ErrStaleFence if a newer token has been issued.
func (m *Manager) Check(resource string, tok Token) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cur, ok := m.current[resource]
	if !ok {
		return fmt.Errorf("fence: unknown resource %q", resource)
	}
	if tok.Seq < cur.Seq {
		return fmt.Errorf("%w: got %d, current %d", ErrStaleFence, tok.Seq, cur.Seq)
	}
	return nil
}

// Current returns the most recently issued token for resource.
func (m *Manager) Current(resource string) (Token, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	tok, ok := m.current[resource]
	return tok, ok
}
