package deadline_test

import (
	"errors"
	"testing"
	"time"

	"github.com/vaultop/internal/deadline"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func newTracker(now time.Time) *deadline.Tracker {
	tr := deadline.New()
	// Swap the internal clock via the exported constructor — we patch via
	// a thin helper that re-uses the same zero-value trick used elsewhere.
	_ = now // used in sub-tests via direct time arithmetic
	return tr
}

func TestStart_And_Check_WithinDeadline(t *testing.T) {
	tr := deadline.New()
	tr.Start("key1", 10*time.Second)
	if err := tr.Check("key1"); err != nil {
		t.Fatalf("expected no error within deadline, got %v", err)
	}
}

func TestCheck_UnknownKey_ReturnsError(t *testing.T) {
	tr := deadline.New()
	err := tr.Check("missing")
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestFinish_RemovesKey(t *testing.T) {
	tr := deadline.New()
	tr.Start("key1", 10*time.Second)
	tr.Finish("key1")
	if err := tr.Check("key1"); err == nil {
		t.Fatal("expected error after Finish, got nil")
	}
}

func TestFinish_Idempotent(t *testing.T) {
	tr := deadline.New()
	tr.Finish("nonexistent") // must not panic
}

func TestEntry_Exceeded_And_Remaining(t *testing.T) {
	now := time.Now()
	e := deadline.Entry{
		Key:       "x",
		Deadline:  now.Add(5 * time.Second),
		StartedAt: now,
	}
	if e.Exceeded(now) {
		t.Error("should not be exceeded yet")
	}
	if e.Exceeded(now.Add(6 * time.Second)) {
		// expected
	} else {
		t.Error("should be exceeded after deadline")
	}
	if r := e.Remaining(now); r == 0 {
		t.Error("remaining should be non-zero")
	}
	if r := e.Remaining(now.Add(10 * time.Second)); r != 0 {
		t.Errorf("remaining should be 0 after expiry, got %v", r)
	}
}

func TestActive_ReturnsOnlyPending(t *testing.T) {
	tr := deadline.New()
	tr.Start("fast", 1*time.Hour)
	tr.Start("slow", 1*time.Hour)
	active := tr.Active()
	if len(active) != 2 {
		t.Fatalf("expected 2 active, got %d", len(active))
	}
}

func TestViolations_ReturnsExpired(t *testing.T) {
	tr := deadline.New()
	tr.Start("ok", 1*time.Hour)
	// Manually craft an expired entry by starting with a negative TTL.
	// Since we cannot inject a clock without exporting it, we start with
	// a very small TTL and sleep briefly.
	tr.Start("expired", 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)

	v := tr.Violations()
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "expired" {
		t.Errorf("unexpected violation key: %s", v[0].Key)
	}
}

func TestCheck_ExpiredKey_ReturnsErrExceeded(t *testing.T) {
	tr := deadline.New()
	tr.Start("k", 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	err := tr.Check("k")
	if !errors.Is(err, deadline.ErrExceeded) {
		t.Fatalf("expected ErrExceeded, got %v", err)
	}
}
