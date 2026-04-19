package ratelimit

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_WithinLimit(t *testing.T) {
	l := New(3, time.Minute)
	for i := 0; i < 3; i++ {
		if err := l.Allow("key1"); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	l := New(2, time.Minute)
	_ = l.Allow("key1")
	_ = l.Allow("key1")
	if err := l.Allow("key1"); err != ErrRateLimited {
		t.Fatalf("expected ErrRateLimited, got %v", err)
	}
}

func TestAllow_WindowExpiry(t *testing.T) {
	base := time.Now()
	l := New(1, time.Minute)
	l.now = fixedClock(base)
	_ = l.Allow("key1")

	// advance past the window
	l.now = fixedClock(base.Add(61 * time.Second))
	if err := l.Allow("key1"); err != nil {
		t.Fatalf("expected nil after window expiry, got %v", err)
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	l := New(1, time.Minute)
	_ = l.Allow("a")
	if err := l.Allow("b"); err != nil {
		t.Fatalf("keys should be independent, got %v", err)
	}
}

func TestReset_ClearsHistory(t *testing.T) {
	l := New(1, time.Minute)
	_ = l.Allow("key1")
	l.Reset("key1")
	if err := l.Allow("key1"); err != nil {
		t.Fatalf("expected nil after reset, got %v", err)
	}
}

func TestRemaining_DecreasesOnAllow(t *testing.T) {
	l := New(3, time.Minute)
	if r := l.Remaining("k"); r != 3 {
		t.Fatalf("expected 3, got %d", r)
	}
	_ = l.Allow("k")
	if r := l.Remaining("k"); r != 2 {
		t.Fatalf("expected 2, got %d", r)
	}
}

func TestRemaining_NeverNegative(t *testing.T) {
	l := New(1, time.Minute)
	_ = l.Allow("k")
	_ = l.Allow("k")
	if r := l.Remaining("k"); r != 0 {
		t.Fatalf("expected 0, got %d", r)
	}
}
