package schedule_test

import (
	"testing"
	"time"

	"github.com/vaultop/vaultop/internal/schedule"
)

func TestCheckAll_MixedDue(t *testing.T) {
	schedules := []schedule.SecretSchedule{
		{Key: "db/password", Policy: schedule.Policy{IntervalDays: 7, LastRotated: time.Now().UTC().AddDate(0, 0, -8)}},
		{Key: "api/key", Policy: schedule.Policy{IntervalDays: 30, LastRotated: time.Now().UTC().AddDate(0, 0, -2)}},
	}

	results, err := schedule.CheckAll(schedules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Due {
		t.Error("expected db/password to be due")
	}
	if results[1].Due {
		t.Error("expected api/key to not be due")
	}
}

func TestCheckAll_InvalidPolicy(t *testing.T) {
	schedules := []schedule.SecretSchedule{
		{Key: "bad/secret", Policy: schedule.Policy{IntervalDays: -5}},
	}
	_, err := schedule.CheckAll(schedules)
	if err == nil {
		t.Fatal("expected error for invalid policy")
	}
}

func TestFilterDue_ReturnsOnlyDue(t *testing.T) {
	results := []schedule.DueResult{
		{Key: "a", Due: true},
		{Key: "b", Due: false},
		{Key: "c", Due: true},
	}
	due := schedule.FilterDue(results)
	if len(due) != 2 {
		t.Fatalf("expected 2 due results, got %d", len(due))
	}
	for _, r := range due {
		if !r.Due {
			t.Errorf("unexpected non-due result: %s", r.Key)
		}
	}
}

func TestFilterDue_Nonedue(t *testing.T) {
	results := []schedule.DueResult{
		{Key: "x", Due: false},
	}
	due := schedule.FilterDue(results)
	if len(due) != 0 {
		t.Fatalf("expected 0 due results, got %d", len(due))
	}
}
