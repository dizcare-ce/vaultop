package lock

import (
	"fmt"
	"os"
	"path/filepath"
)

// Manager manages named locks under a shared directory.
type Manager struct {
	dir string
}

// NewManager returns a Manager that stores lock files under dir.
func NewManager(dir string) (*Manager, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("lock manager: mkdir: %w", err)
	}
	return &Manager{dir: dir}, nil
}

// Acquire acquires a named lock. Returns ErrLocked if already held.
func (m *Manager) Acquire(name string) (*Lock, error) {
	l := New(m.lockPath(name))
	if err := l.Acquire(); err != nil {
		return nil, err
	}
	return l, nil
}

// Release releases a named lock.
func (m *Manager) Release(name string) error {
	return New(m.lockPath(name)).Release()
}

// IsHeld reports whether a named lock is currently held.
func (m *Manager) IsHeld(name string) bool {
	return New(m.lockPath(name)).IsHeld()
}

func (m *Manager) lockPath(name string) string {
	return filepath.Join(m.dir, name+".lock")
}
