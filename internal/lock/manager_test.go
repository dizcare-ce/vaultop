package lock_test

import (
	"errors"
	"testing"

	"github.com/yourusername/vaultop/internal/lock"
)

func TestManager_AcquireAndRelease(t *testing.T) {
	m, err := lock.NewManager(t.TempDir())
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	l, err := m.Acquire("rotation")
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	if !m.IsHeld("rotation") {
		t.Fatal("expected lock to be held")
	}
	l.Release()
	if m.IsHeld("rotation") {
		t.Fatal("expected lock to be released")
	}
}

func TestManager_DoubleAcquire_ReturnsErrLocked(t *testing.T) {
	m, err := lock.NewManager(t.TempDir())
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	l, err := m.Acquire("rotation")
	if err != nil {
		t.Fatalf("first Acquire: %v", err)
	}
	defer l.Release()

	_, err = m.Acquire("rotation")
	if !errors.Is(err, lock.ErrLocked) {
		t.Fatalf("expected ErrLocked, got: %v", err)
	}
}

func TestManager_IndependentNames(t *testing.T) {
	m, err := lock.NewManager(t.TempDir())
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	l1, err := m.Acquire("alpha")
	if err != nil {
		t.Fatalf("Acquire alpha: %v", err)
	}
	defer l1.Release()

	l2, err := m.Acquire("beta")
	if err != nil {
		t.Fatalf("Acquire beta: %v", err)
	}
	defer l2.Release()
}

func TestManager_ReleaseAllowsReacquire(t *testing.T) {
	m, err := lock.NewManager(t.TempDir())
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	l, err := m.Acquire("rotation")
	if err != nil {
		t.Fatalf("first Acquire: %v", err)
	}
	l.Release()

	l2, err := m.Acquire("rotation")
	if err != nil {
		t.Fatalf("re-Acquire after Release: %v", err)
	}
	defer l2.Release()

	if !m.IsHeld("rotation") {
		t.Fatal("expected lock to be held after re-acquire")
	}
}
