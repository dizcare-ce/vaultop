// Package alert provides threshold-based alerting for secret metrics and TTL expiry.
package alert

import (
	"fmt"
	"io"
	"time"
)

// Level represents alert severity.
type Level string

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelCrit  Level = "crit"
)

// Alert holds a single alert event.
type Alert struct {
	Key       string
	Level     Level
	Message   string
	Timestamp time.Time
}

// Rule defines when an alert should fire.
type Rule struct {
	Key            string
	WarnWithin     time.Duration // warn if expiry is within this window
	CritWithin     time.Duration // crit if expiry is within this window
}

// Notifier sends alerts somewhere.
type Notifier interface {
	Notify(a Alert) error
}

// WriterNotifier writes alerts as text lines to a writer.
type WriterNotifier struct {
	w io.Writer
}

func NewWriterNotifier(w io.Writer) *WriterNotifier {
	return &WriterNotifier{w: w}
}

func (wn *WriterNotifier) Notify(a Alert) error {
	_, err := fmt.Fprintf(wn.w, "[%s] %s %s: %s\n",
		a.Timestamp.Format(time.RFC3339), a.Level, a.Key, a.Message)
	return err
}

// Noop discards all alerts.
type Noop struct{}

func (Noop) Notify(Alert) error { return nil }

// CheckExpiry fires an alert if the given expiry time breaches rule thresholds.
func CheckExpiry(rule Rule, expiry time.Time, now time.Time) (Alert, bool) {
	remaining := expiry.Sub(now)
	switch {
	case remaining <= rule.CritWithin:
		return Alert{Key: rule.Key, Level: LevelCrit,
			Message: fmt.Sprintf("expires in %s", remaining.Round(time.Second)),
			Timestamp: now}, true
	case remaining <= rule.WarnWithin:
		return Alert{Key: rule.Key, Level: LevelWarn,
			Message: fmt.Sprintf("expires in %s", remaining.Round(time.Second)),
			Timestamp: now}, true
	}
	return Alert{}, false
}
