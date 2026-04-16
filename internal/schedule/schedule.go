// Package schedule provides rotation scheduling based on cron expressions
// and interval-based policies defined in the vaultop configuration.
package schedule

import (
	"fmt"
	"time"
)

// Policy describes when a secret should be rotated.
type Policy struct {
	// IntervalDays is the number of days between rotations.
	IntervalDays int `yaml:"interval_days"`
	// LastRotated is the timestamp of the most recent rotation.
	LastRotated time.Time `yaml:"last_rotated"`
}

// IsDue reports whether the policy indicates a rotation is due as of now.
func (p Policy) IsDue() bool {
	return p.IsDueAt(time.Now().UTC())
}

// IsDueAt reports whether a rotation is due at the given time t.
func (p Policy) IsDueAt(t time.Time) bool {
	if p.IntervalDays <= 0 {
		return false
	}
	if p.LastRotated.IsZero() {
		return true
	}
	next := p.LastRotated.AddDate(0, 0, p.IntervalDays)
	return !t.Before(next)
}

// NextRotation returns the time at which the next rotation is due.
func (p Policy) NextRotation() (time.Time, error) {
	if p.IntervalDays <= 0 {
		return time.Time{}, fmt.Errorf("schedule: interval_days must be greater than 0")
	}
	if p.LastRotated.IsZero() {
		return time.Now().UTC(), nil
	}
	return p.LastRotated.AddDate(0, 0, p.IntervalDays), nil
}

// Validate returns an error if the policy is misconfigured.
func (p Policy) Validate() error {
	if p.IntervalDays < 0 {
		return fmt.Errorf("schedule: interval_days cannot be negative")
	}
	return nil
}
