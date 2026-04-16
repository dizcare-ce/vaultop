package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/vaultop/vaultop/internal/audit"
)

func TestLog_WritesJSONEntry(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	err := l.Log(audit.EventRotated, "stub", "db/password", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	var entry audit.Entry
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		t.Fatalf("failed to unmarshal entry: %v", err)
	}

	if entry.Event != audit.EventRotated {
		t.Errorf("expected event %q, got %q", audit.EventRotated, entry.Event)
	}
	if entry.Provider != "stub" {
		t.Errorf("expected provider %q, got %q", "stub", entry.Provider)
	}
	if entry.SecretKey != "db/password" {
		t.Errorf("expected secret_key %q, got %q", "db/password", entry.SecretKey)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLog_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	events := []audit.EventType{audit.EventRotated, audit.EventDryRun, audit.EventFailed}
	for _, ev := range events {
		if err := l.Log(ev, "stub", "key", "msg"); err != nil {
			t.Fatalf("Log error: %v", err)
		}
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}

func TestNewFileLogger_InvalidPath(t *testing.T) {
	_, err := audit.NewFileLogger("/nonexistent_dir/audit.log")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}
