package lock_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultop/internal/lock"
)

func tempLockPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "vaultop.lock")
}

func TestAcquire_CreatesFile(t *testing.T) {
	path := tempLockPath(t)
	l := lock.New(path)
	if err := l.Acquire(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer l.Release()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("lock file not created: %v", err)
	}
}

func TestAcquire_FailsWhenAlreadyHeld(t *testing.T) {
	path := tempLockPath(t)
	l := lock.New(path)
	if err := l.Acquire(); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	defer l.Release()

	l2 := lock.New(path)
	err := l2.Acquire()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, lock.ErrLocked) {
		t.Fatalf("expected ErrLocked, got: %v", err)
	}
}

func TestRelease_RemovesFile(t *testing.T) {
	path := tempLockPath(t)
	l := lock.New(path)
	_ = l.Acquire()
	if err := l.Release(); err != nil {
		t.Fatalf("release failed: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("lock file still exists after release")
	}
}

func TestRelease_IdempotentWhenNotHeld(t *testing.T) {
	path := tempLockPath(t)
	l := lock.New(path)
	if err := l.Release(); err != nil {
		t.Fatalf("unexpected error on release of non-existent lock: %v", err)
	}
}

func TestIsHeld(t *testing.T) {
	path := tempLockPath(t)
	l := lock.New(path)
	if l.IsHeld() {
		t.Fatal("expected lock to not be held")
	}
	_ = l.Acquire()
	if !l.IsHeld() {
		t.Fatal("expected lock to be held")
	}
	_ = l.Release()
	if l.IsHeld() {
		t.Fatal("expected lock to not be held after release")
	}
}
