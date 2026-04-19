package backup_test

import (
	"os"
	"testing"

	"github.com/vaultop/internal/backup"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "backup-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func newManager(t *testing.T) *backup.Manager {
	t.Helper()
	m, err := backup.NewManager(tempDir(t))
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	return m
}

func TestSave_And_Load_RoundTrip(t *testing.T) {
	m := newManager(t)
	secrets := map[string]string{"db/pass": "s3cr3t", "api/key": "abc123"}
	id, err := m.Save(secrets)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty id")
	}
	entry, err := m.Load(id)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if entry.ID != id {
		t.Errorf("id mismatch: got %s want %s", entry.ID, id)
	}
	for k, v := range secrets {
		if entry.Secrets[k] != v {
			t.Errorf("key %s: got %q want %q", k, entry.Secrets[k], v)
		}
	}
}

func TestLoad_UnknownID_ReturnsError(t *testing.T) {
	m := newManager(t)
	_, err := m.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown id")
	}
}

func TestList_ReturnsAllIDs(t *testing.T) {
	m := newManager(t)
	for i := 0; i < 3; i++ {
		if _, err := m.Save(map[string]string{"k": "v"}); err != nil {
			t.Fatalf("Save: %v", err)
		}
	}
	ids, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(ids) != 3 {
		t.Errorf("expected 3 entries, got %d", len(ids))
	}
}

func TestList_EmptyDir_ReturnsEmptySlice(t *testing.T) {
	m := newManager(t)
	ids, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected 0 entries, got %d", len(ids))
	}
}
