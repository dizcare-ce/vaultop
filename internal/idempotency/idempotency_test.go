package idempotency

import (
	"errors"
	"testing"
	"time"
)

var (
	t0     = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	window = 5 * time.Minute
)

func fixedClock(at time.Time) func() time.Time {
	return func() time.Time { return at }
}

func TestCheck_NewKey_ReturnsFalse(t *testing.T) {
	s := newWithClock(window, fixedClock(t0))
	_, ok := s.Check("op-1")
	if ok {
		t.Fatal("expected false for unseen key")
	}
}

func TestRecord_And_Check_ReturnsResult(t *testing.T) {
	s := newWithClock(window, fixedClock(t0))
	r := Result{Value: "abc123", Err: nil, CompletedAt: t0}
	s.Record("op-1", r)

	got, ok := s.Check("op-1")
	if !ok {
		t.Fatal("expected true for recorded key")
	}
	if got.Value != "abc123" {
		t.Fatalf("expected abc123, got %s", got.Value)
	}
}

func TestCheck_ExpiredKey_ReturnsFalse(t *testing.T) {
	s := newWithClock(window, fixedClock(t0))
	s.Record("op-1", Result{Value: "v"})

	// advance clock beyond window
	s.clock = fixedClock(t0.Add(window + time.Second))
	_, ok := s.Check("op-1")
	if ok {
		t.Fatal("expected false for expired key")
	}
}

func TestRecord_StoresError(t *testing.T) {
	s := newWithClock(window, fixedClock(t0))
	sentinel := errors.New("provider unavailable")
	s.Record("op-err", Result{Err: sentinel})

	got, ok := s.Check("op-err")
	if !ok {
		t.Fatal("expected recorded entry")
	}
	if !errors.Is(got.Err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", got.Err)
	}
}

func TestPurge_RemovesExpiredEntries(t *testing.T) {
	s := newWithClock(window, fixedClock(t0))
	s.Record("op-a", Result{Value: "a"})
	s.Record("op-b", Result{Value: "b"})

	s.clock = fixedClock(t0.Add(window + time.Second))
	s.Purge()

	if s.Len() != 0 {
		t.Fatalf("expected 0 entries after purge, got %d", s.Len())
	}
}

func TestLen_CountsOnlyActive(t *testing.T) {
	s := newWithClock(window, fixedClock(t0))
	s.Record("op-a", Result{Value: "a"})
	s.Record("op-b", Result{Value: "b"})

	if s.Len() != 2 {
		t.Fatalf("expected 2, got %d", s.Len())
	}

	s.clock = fixedClock(t0.Add(window + time.Second))
	if s.Len() != 0 {
		t.Fatalf("expected 0 after expiry, got %d", s.Len())
	}
}

func TestRecord_OverwritesPreviousEntry(t *testing.T) {
	s := newWithClock(window, fixedClock(t0))
	s.Record("op-1", Result{Value: "first"})
	s.Record("op-1", Result{Value: "second"})

	got, ok := s.Check("op-1")
	if !ok {
		t.Fatal("expected recorded entry")
	}
	if got.Value != "second" {
		t.Fatalf("expected second, got %s", got.Value)
	}
}
