package history_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultop/internal/history"
)

func TestReport_AllEntries(t *testing.T) {
	s, _ := history.Load(tempPath(t))
	now := time.Now().UTC()
	_ = s.Record(history.Entry{SecretKey: "a/key", Provider: "aws", RotatedAt: now, Success: true})
	_ = s.Record(history.Entry{SecretKey: "b/key", Provider: "gcp", RotatedAt: now, Success: false})

	var buf bytes.Buffer
	if err := history.Report(&buf, s, nil, time.Time{}); err != nil {
		t.Fatalf("report: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "a/key") {
		t.Error("expected a/key in output")
	}
	if !strings.Contains(out, "FAILED") {
		t.Error("expected FAILED status in output")
	}
}

func TestReport_FilterByKey(t *testing.T) {
	s, _ := history.Load(tempPath(t))
	now := time.Now().UTC()
	_ = s.Record(history.Entry{SecretKey: "a/key", Provider: "aws", RotatedAt: now, Success: true})
	_ = s.Record(history.Entry{SecretKey: "b/key", Provider: "gcp", RotatedAt: now, Success: true})

	var buf bytes.Buffer
	_ = history.Report(&buf, s, []string{"a/key"}, time.Time{})
	out := buf.String()
	if strings.Contains(out, "b/key") {
		t.Error("b/key should be filtered out")
	}
}

func TestReport_FilterBySince(t *testing.T) {
	s, _ := history.Load(tempPath(t))
	now := time.Now().UTC()
	_ = s.Record(history.Entry{SecretKey: "old", Provider: "aws", RotatedAt: now.Add(-72 * time.Hour), Success: true})
	_ = s.Record(history.Entry{SecretKey: "new", Provider: "aws", RotatedAt: now, Success: true})

	var buf bytes.Buffer
	_ = history.Report(&buf, s, nil, now.Add(-24*time.Hour))
	out := buf.String()
	if strings.Contains(out, "old") {
		t.Error("old entry should be filtered by since")
	}
	if !strings.Contains(out, "new") {
		t.Error("new entry should appear")
	}
}
