package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultop/internal/history"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestLoad_NewFile(t *testing.T) {
	s, err := history.Load(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.All()) != 0 {
		t.Fatalf("expected empty store")
	}
}

func TestRecord_PersistsAndReloads(t *testing.T) {
	p := tempPath(t)
	s, _ := history.Load(p)
	e := history.Entry{
		SecretKey: "db/password",
		RotatedAt: time.Now().UTC().Truncate(time.Second),
		Provider:  "aws",
		Success:   true,
	}
	if err := s.Record(e); err != nil {
		t.Fatalf("record: %v", err)
	}
	s2, err := history.Load(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(s2.All()) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(s2.All()))
	}
	if s2.All()[0].SecretKey != "db/password" {
		t.Errorf("unexpected key: %s", s2.All()[0].SecretKey)
	}
}

func TestLastRotated_ReturnsLatestSuccess(t *testing.T) {
	s, _ := history.Load(tempPath(t))
	now := time.Now().UTC()
	_ = s.Record(history.Entry{SecretKey: "k", RotatedAt: now.Add(-48 * time.Hour), Provider: "stub", Success: true})
	_ = s.Record(history.Entry{SecretKey: "k", RotatedAt: now.Add(-24 * time.Hour), Provider: "stub", Success: false})
	_ = s.Record(history.Entry{SecretKey: "k", RotatedAt: now.Add(-1 * time.Hour), Provider: "stub", Success: true})

	got := s.LastRotated("k")
	if got.IsZero() {
		t.Fatal("expected non-zero time")
	}
	if got.Before(now.Add(-2 * time.Hour)) {
		t.Errorf("expected most recent success, got %v", got)
	}
}

func TestLastRotated_UnknownKey(t *testing.T) {
	s, _ := history.Load(tempPath(t))
	if !s.LastRotated("missing").IsZero() {
		t.Error("expected zero time for unknown key")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not json"), 0600)
	_, err := history.Load(p)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
