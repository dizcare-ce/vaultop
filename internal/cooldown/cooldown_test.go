package cooldown

import (
	"testing"
	"time"
)

// fixedClock returns a clock whose value can be advanced manually.
func fixedClock(initial time.Time) (Clock, func(d time.Duration)) {
	current := initial
	clock := func() time.Time { return current }
	advance := func(d time.Duration) { current = current.Add(d) }
	return clock, advance
}

func TestAllow_FirstCall_Succeeds(t *testing.T) {
	clock, _ := fixedClock(time.Now())
	m := newWithClock(5*time.Second, clock)

	if err := m.Allow("key1"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAllow_WithinCooldown_ReturnsErrCooldown(t *testing.T) {
	clock, advance := fixedClock(time.Now())
	m := newWithClock(10*time.Second, clock)

	_ = m.Allow("key1")
	advance(3 * time.Second)

	if err := m.Allow("key1"); err != ErrCooldown {
		t.Fatalf("expected ErrCooldown, got %v", err)
	}
}

func TestAllow_AfterCooldown_Succeeds(t *testing.T) {
	clock, advance := fixedClock(time.Now())
	m := newWithClock(5*time.Second, clock)

	_ = m.Allow("key1")
	advance(6 * time.Second)

	if err := m.Allow("key1"); err != nil {
		t.Fatalf("expected nil after cooldown elapsed, got %v", err)
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	clock, _ := fixedClock(time.Now())
	m := newWithClock(10*time.Second, clock)

	_ = m.Allow("a")

	if err := m.Allow("b"); err != nil {
		t.Fatalf("key 'b' should not be affected by key 'a' cooldown, got %v", err)
	}
}

func TestReset_AllowsImmediateReuse(t *testing.T) {
	clock, _ := fixedClock(time.Now())
	m := newWithClock(30*time.Second, clock)

	_ = m.Allow("key1")
	m.Reset("key1")

	if err := m.Allow("key1"); err != nil {
		t.Fatalf("expected nil after reset, got %v", err)
	}
}

func TestRemaining_ZeroWhenNotCooling(t *testing.T) {
	clock, _ := fixedClock(time.Now())
	m := newWithClock(10*time.Second, clock)

	if r := m.Remaining("unknown"); r != 0 {
		t.Fatalf("expected 0 for unknown key, got %v", r)
	}
}

func TestRemaining_ReturnsPositiveDuringCooldown(t *testing.T) {
	clock, advance := fixedClock(time.Now())
	m := newWithClock(10*time.Second, clock)

	_ = m.Allow("key1")
	advance(3 * time.Second)

	r := m.Remaining("key1")
	if r <= 0 {
		t.Fatalf("expected positive remaining, got %v", r)
	}
	if r > 7*time.Second+time.Millisecond {
		t.Fatalf("remaining %v exceeds expected ~7s", r)
	}
}
