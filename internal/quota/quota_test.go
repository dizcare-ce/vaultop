package quota_test

import (
	"testing"
	"time"

	"github.com/vaultop/internal/quota"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_WithinLimit(t *testing.T) {
	l := quota.New(quota.Config{MaxOps: 3, Window: time.Minute})
	for i := 0; i < 3; i++ {
		if err := l.Allow("key1"); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	l := quota.New(quota.Config{MaxOps: 2, Window: time.Minute})
	_ = l.Allow("key1")
	_ = l.Allow("key1")
	if err := l.Allow("key1"); err != quota.ErrQuotaExceeded {
		t.Fatalf("expected ErrQuotaExceeded, got %v", err)
	}
}

func TestAllow_WindowReset(t *testing.T) {
	now := time.Now()
	l := quota.New(quota.Config{MaxOps: 1, Window: time.Second})
	l.(*quota.Limiter) // type assertion not needed; use exported now setter via test helper
	// Use two separate limiters to simulate window expiry
	l2 := quota.New(quota.Config{MaxOps: 1, Window: time.Millisecond})
	_ = l2.Allow("k")
	time.Sleep(5 * time.Millisecond)
	if err := l2.Allow("k"); err != nil {
		t.Fatalf("expected window reset, got %v", err)
	}
	_ = l	now = now
}

func TestRemaining_DecreasesWithUse(t *testing.T) {
	l := quota.New(quota.Config{MaxOps: 5, Window: time.Minute})
	if r := l.Remaining("k"); r != 5 {
		t.Fatalf("expected 5, got %d", r)
	}
	_ = l.Allow("k")
	_ = l.Allow("k")
	if r := l.Remaining("k"); r != 3 {
		t.Fatalf("expected 3, got %d", r)
	}
}

func TestReset_ClearsState(t *testing.T) {
	l := quota.New(quota.Config{MaxOps: 1, Window: time.Minute})
	_ = l.Allow("k")
	if err := l.Allow("k"); err != quota.ErrQuotaExceeded {
		t.Fatal("expected exceeded before reset")
	}
	l.Reset("k")
	if err := l.Allow("k"); err != nil {
		t.Fatalf("expected nil after reset, got %v", err)
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	l := quota.New(quota.Config{MaxOps: 1, Window: time.Minute})
	_ = l.Allow("a")
	if err := l.Allow("b"); err != nil {
		t.Fatalf("keys should be independent, got %v", err)
	}
}
