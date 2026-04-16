package schedule_test

import (
	"testing"
	"time"

	"github.com/vaultop/vaultop/internal/schedule"
)

func daysAgo(n int) time.Time {
	return time.Now().UTC().AddDate(0, 0, -n)
}

func TestIsDueAt_NeverRotated(t *testing.T) {
	p := schedule.Policy{IntervalDays: 7}
	if !p.IsDue() {
		t.Fatal("expected policy with zero LastRotated to be due")
	}
}

func TestIsDueAt_NotYetDue(t *testing.T) {
	p := schedule.Policy{IntervalDays: 7, LastRotated: daysAgo(3)}
	if p.IsDue() {
		t.Fatal("expected policy rotated 3 days ago with 7-day interval to not be due")
	}
}

func TestIsDueAt_ExactlyDue(t *testing.T) {
	last := daysAgo(7)
	p := schedule.Policy{IntervalDays: 7, LastRotated: last}
	check := last.AddDate(0, 0, 7)
	if !p.IsDueAt(check) {
		t.Fatal("expected policy to be due exactly on rotation day")
	}
}

func TestIsDueAt_ZeroInterval(t *testing.T) {
	p := schedule.Policy{IntervalDays: 0}
	if p.IsDue() {
		t.Fatal("expected zero interval to never be due")
	}
}

func TestNextRotation_NoLastRotated(t *testing.T) {
	p := schedule.Policy{IntervalDays: 30}
	next, err := p.NextRotation()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if next.IsZero() {
		t.Fatal("expected non-zero next rotation time")
	}
}

func TestNextRotation_InvalidInterval(t *testing.T) {
	p := schedule.Policy{IntervalDays: 0}
	_, err := p.NextRotation()
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestValidate_Negative(t *testing.T) {
	p := schedule.Policy{IntervalDays: -1}
	if err := p.Validate(); err == nil {
		t.Fatal("expected validation error for negative interval")
	}
}

func TestValidate_Valid(t *testing.T) {
	p := schedule.Policy{IntervalDays: 14}
	if err := p.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}
