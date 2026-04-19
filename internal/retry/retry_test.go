package retry_test

import (
	"errors"
	"testing"
	"time"

	"github.com/vaultop/internal/retry"
)

var errTemp = errors.New("temporary error")

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(retry.DefaultConfig(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesAndSucceeds(t *testing.T) {
	calls := 0
	cfg := retry.Config{MaxAttempts: 3, Delay: time.Millisecond, Multiplier: 1.0}
	err := retry.Do(cfg, func() error {
		calls++
		if calls < 3 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	calls := 0
	cfg := retry.Config{MaxAttempts: 3, Delay: time.Millisecond, Multiplier: 1.0}
	err := retry.Do(cfg, func() error {
		calls++
		return errTemp
	})
	if !errors.Is(err, retry.ErrMaxAttemptsReached) {
		t.Fatalf("expected ErrMaxAttemptsReached, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ZeroMaxAttempts_RunsOnce(t *testing.T) {
	calls := 0
	cfg := retry.Config{MaxAttempts: 0, Delay: time.Millisecond, Multiplier: 1.0}
	retry.Do(cfg, func() error { calls++; return errTemp }) //nolint
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}
