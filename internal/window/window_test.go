package window

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRecord_And_Count(t *testing.T) {
	now := time.Now()
	c := newWithClock(10*time.Second, fixedClock(now))
	c.Record("k")
	c.Record("k")
	if got := c.Count("k"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestCount_Zero_ForUnknownKey(t *testing.T) {
	c := New(10 * time.Second)
	if got := c.Count("missing"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestEviction_RemovesOldEntries(t *testing.T) {
	base := time.Now()
	var now time.Time
	c := newWithClock(5*time.Second, func() time.Time { return now })

	// record at t=0
	now = base
	c.Record("k")
	c.Record("k")

	// advance past window
	now = base.Add(6 * time.Second)
	c.Record("k") // one new event

	if got := c.Count("k"); got != 1 {
		t.Fatalf("expected 1 after eviction, got %d", got)
	}
}

func TestReset_ClearsKey(t *testing.T) {
	now := time.Now()
	c := newWithClock(10*time.Second, fixedClock(now))
	c.Record("k")
	c.Record("k")
	c.Reset("k")
	if got := c.Count("k"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestIndependentKeys(t *testing.T) {
	now := time.Now()
	c := newWithClock(10*time.Second, fixedClock(now))
	c.Record("a")
	c.Record("a")
	c.Record("b")
	if got := c.Count("a"); got != 2 {
		t.Fatalf("key a: expected 2, got %d", got)
	}
	if got := c.Count("b"); got != 1 {
		t.Fatalf("key b: expected 1, got %d", got)
	}
}

func TestBoundaryEvent_IncludedInWindow(t *testing.T) {
	base := time.Now()
	var now time.Time
	c := newWithClock(5*time.Second, func() time.Time { return now })

	// record exactly at the boundary (cutoff == entry.at, should be kept)
	now = base
	c.Record("k")

	now = base.Add(5 * time.Second)
	if got := c.Count("k"); got != 1 {
		t.Fatalf("boundary event should be included, got %d", got)
	}
}
