package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var baseRule = Rule{
	Key:        "db/password",
	WarnWithin: 7 * 24 * time.Hour,
	CritWithin: 24 * time.Hour,
}

func TestCheckExpiry_Crit(t *testing.T) {
	now := time.Now()
	expiry := now.Add(12 * time.Hour)
	a, fired := CheckExpiry(baseRule, expiry, now)
	if !fired {
		t.Fatal("expected alert to fire")
	}
	if a.Level != LevelCrit {
		t.Fatalf("expected crit, got %s", a.Level)
	}
}

func TestCheckExpiry_Warn(t *testing.T) {
	now := time.Now()
	expiry := now.Add(3 * 24 * time.Hour)
	a, fired := CheckExpiry(baseRule, expiry, now)
	if !fired {
		t.Fatal("expected alert to fire")
	}
	if a.Level != LevelWarn {
		t.Fatalf("expected warn, got %s", a.Level)
	}
}

func TestCheckExpiry_NoAlert(t *testing.T) {
	now := time.Now()
	expiry := now.Add(30 * 24 * time.Hour)
	_, fired := CheckExpiry(baseRule, expiry, now)
	if fired {
		t.Fatal("expected no alert")
	}
}

func TestWriterNotifier_Notify(t *testing.T) {
	var buf bytes.Buffer
	n := NewWriterNotifier(&buf)
	a := Alert{Key: "k", Level: LevelWarn, Message: "soon", Timestamp: time.Now()}
	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "warn") {
		t.Fatalf("expected 'warn' in output, got: %s", buf.String())
	}
}

func TestNoop_Notify(t *testing.T) {
	var n Noop
	if err := n.Notify(Alert{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
