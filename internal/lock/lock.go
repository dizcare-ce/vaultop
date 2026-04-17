// Package lock provides a simple file-based locking mechanism to prevent
// concurrent rotations from running simultaneously.
package lock

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ErrLocked is returned when the lock is already held.
var ErrLocked = errors.New("lock: already locked")

// Lock represents a file-based lock.
type Lock struct {
	path string
}

// New returns a Lock backed by the given file path.
func New(path string) *Lock {
	return &Lock{path: filepath.Clean(path)}
}

// Acquire attempts to create the lock file. Returns ErrLocked if already held.
func (l *Lock) Acquire() error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		if os.IsExist(err) {
			owner, _ := l.readOwner()
			return fmt.Errorf("%w (held by pid %s)", ErrLocked, owner)
		}
		return fmt.Errorf("lock: create: %w", err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%d\n%d", os.Getpid(), time.Now().Unix())
	return err
}

// Release removes the lock file.
func (l *Lock) Release() error {
	if err := os.Remove(l.path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("lock: release: %w", err)
	}
	return nil
}

// IsHeld reports whether the lock file currently exists.
func (l *Lock) IsHeld() bool {
	_, err := os.Stat(l.path)
	return err == nil
}

func (l *Lock) readOwner() (string, error) {
	data, err := os.ReadFile(l.path)
	if err != nil {
		return "unknown", err
	}
	parts := strings.SplitN(string(data), "\n", 2)
	pid, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return "unknown", nil
	}
	return strconv.Itoa(pid), nil
}
