package schedule

import (
	"fmt"
	"time"
)

// SecretSchedule pairs a secret key with its rotation policy.
type SecretSchedule struct {
	Key    string
	Policy Policy
}

// DueResult holds the outcome of a due-check for a single secret.
type DueResult struct {
	Key         string
	Due         bool
	NextRotation time.Time
}

// CheckAll evaluates each SecretSchedule and returns a DueResult per entry.
func CheckAll(schedules []SecretSchedule) ([]DueResult, error) {
	results := make([]DueResult, 0, len(schedules))
	for _, s := range schedules {
		if err := s.Policy.Validate(); err != nil {
			return nil, fmt.Errorf("schedule: key %q: %w", s.Key, err)
		}
		next, err := s.Policy.NextRotation()
		if err != nil {
			next = time.Time{}
		}
		results = append(results, DueResult{
			Key:          s.Key,
			Due:          s.Policy.IsDue(),
			NextRotation: next,
		})
	}
	return results, nil
}

// FilterDue returns only those DueResults where Due is true.
func FilterDue(results []DueResult) []DueResult {
	var due []DueResult
	for _, r := range results {
		if r.Due {
			due = append(due, r)
		}
	}
	return due
}
