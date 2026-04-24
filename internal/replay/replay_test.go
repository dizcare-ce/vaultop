package replay

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestCheck_FirstTime_Succeeds(t *testing.T) {
	d := New(time.Minute)
	if err := d.Check("op-1"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCheck_DuplicateWithinWindow_ReturnsErrReplay(t *testing.T) {
	d := New(time.Minute)
	_ = d.Check("op-2")
	if err := d.Check("op-2"); err != ErrReplay {
		t.Fatalf("expected ErrReplay, got %v", err)
	}
}

func TestCheck_AfterWindowExpiry_Succeeds(t *testing.T) {
	now := time.Now()
	d := New(time.Minute)
	d.clock = fixedClock(now)

	_ = d.Check("op-3")

	// Advance clock beyond the window.
	d.clock = fixedClock(now.Add(2 * time.Minute))

	if err := d.Check("op-3"); err != nil {
		t.Fatalf("expected nil after expiry, got %v", err)
	}
}

func TestSeen_ReturnsTrueWhenTracked(t *testing.T) {
	d := New(time.Minute)
	_ = d.Check("op-4")
	if !d.Seen("op-4") {
		t.Fatal("expected Seen to return true")
	}
}

func TestSeen_ReturnsFalseWhenNotTracked(t *testing.T) {
	d := New(time.Minute)
	if d.Seen("unknown") {
		t.Fatal("expected Seen to return false for unknown ID")
	}
}

func TestSeen_ReturnsFalseAfterExpiry(t *testing.T) {
	now := time.Now()
	d := New(time.Minute)
	d.clock = fixedClock(now)
	_ = d.Check("op-5")

	d.clock = fixedClock(now.Add(2 * time.Minute))
	if d.Seen("op-5") {
		t.Fatal("expected Seen to return false after expiry")
	}
}

func TestSize_ReflectsActiveEntries(t *testing.T) {
	now := time.Now()
	d := New(time.Minute)
	d.clock = fixedClock(now)

	_ = d.Check("a")
	_ = d.Check("b")

	if s := d.Size(); s != 2 {
		t.Fatalf("expected size 2, got %d", s)
	}

	// Expire entries.
	d.clock = fixedClock(now.Add(2 * time.Minute))
	if s := d.Size(); s != 0 {
		t.Fatalf("expected size 0 after expiry, got %d", s)
	}
}

func TestCheck_IndependentIDs_BothSucceed(t *testing.T) {
	d := New(time.Minute)
	if err := d.Check("x"); err != nil {
		t.Fatalf("unexpected error for x: %v", err)
	}
	if err := d.Check("y"); err != nil {
		t.Fatalf("unexpected error for y: %v", err)
	}
}
