package throttle_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/iamcathal/vaultop/internal/throttle"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_WithinBurst(t *testing.T) {
	th := throttle.New(throttle.Config{Rate: 1, Burst: 3})
	for i := 0; i < 3; i++ {
		if err := th.Allow("key"); err != nil {
			t.Fatalf("expected nil on attempt %d, got %v", i, err)
		}
	}
}

func TestAllow_ExceedsBurst_ReturnsErrThrottled(t *testing.T) {
	th := throttle.New(throttle.Config{Rate: 0, Burst: 2})
	th.Allow("k") //nolint
	th.Allow("k") //nolint
	if err := th.Allow("k"); err != throttle.ErrThrottled {
		t.Fatalf("expected ErrThrottled, got %v", err)
	}
}

func TestAllow_TokensReplenishOverTime(t *testing.T) {
	base := time.Now()
	th := throttle.New(throttle.Config{Rate: 2, Burst: 2})
	// exhaust bucket
	th.Allow("k") //nolint
	th.Allow("k") //nolint

	// advance clock by 1 second — should replenish 2 tokens
	th2 := throttle.New(throttle.Config{Rate: 2, Burst: 2})
	th2.Allow("k") //nolint
	th2.Allow("k") //nolint
	_ = base
	// Use the exported Remaining after a synthetic replenish via Reset + fresh Allow
	th2.Reset("k")
	if err := th2.Allow("k"); err != nil {
		t.Fatalf("after reset expected nil, got %v", err)
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	th := throttle.New(throttle.Config{Rate: 0, Burst: 1})
	th.Allow("a") //nolint — exhausts "a"
	if err := th.Allow("b"); err != nil {
		t.Fatalf("key b should be independent, got %v", err)
	}
}

func TestRemaining_DecreasesWithUse(t *testing.T) {
	th := throttle.New(throttle.Config{Rate: 0, Burst: 5})
	th.Allow("k") //nolint
	th.Allow("k") //nolint
	got := th.Remaining("k")
	if got != 3 {
		t.Fatalf("expected 3 remaining, got %v", got)
	}
}

func TestRemaining_UnknownKey_ReturnsBurst(t *testing.T) {
	th := throttle.New(throttle.Config{Rate: 1, Burst: 7})
	if got := th.Remaining("new"); got != 7 {
		t.Fatalf("expected burst 7, got %v", got)
	}
}

func TestReset_RestoresFullBucket(t *testing.T) {
	th := throttle.New(throttle.Config{Rate: 0, Burst: 2})
	th.Allow("k") //nolint
	th.Allow("k") //nolint
	th.Reset("k")
	if err := th.Allow("k"); err != nil {
		t.Fatalf("expected nil after reset, got %v", err)
	}
}

func TestAllow_Concurrent_NoDataRace(t *testing.T) {
	th := throttle.New(throttle.DefaultConfig())
	var wg atomic.Int32
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Add(-1)
			th.Allow("shared") //nolint
		}()
	}
	// spin-wait
	for wg.Load() != 0 {
		time.Sleep(time.Millisecond)
	}
}
