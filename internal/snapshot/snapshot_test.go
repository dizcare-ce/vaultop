package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultop/internal/provider"
	"github.com/vaultop/internal/snapshot"
)

func newStub(t *testing.T) provider.Provider {
	t.Helper()
	p, err := provider.New(provider.TypeStub, nil)
	if err != nil {
		t.Fatalf("new stub: %v", err)
	}
	return p
}

func TestTake_CapturesSecrets(t *testing.T) {
	p := newStub(t)
	_ = p.Set("alpha", "aaa")
	_ = p.Set("beta", "bbb")

	s, err := snapshot.Take(p, []string{"alpha", "beta"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Secrets["alpha"] != "aaa" || s.Secrets["beta"] != "bbb" {
		t.Errorf("secrets mismatch: %v", s.Secrets)
	}
	if s.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestTake_MissingKey_ReturnsError(t *testing.T) {
	p := newStub(t)
	_, err := snapshot.Take(p, []string{"missing"})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	p := newStub(t)
	_ = p.Set("x", "val")

	s, _ := snapshot.Take(p, []string{"x"})
	path := filepath.Join(t.TempDir(), "snap.json")

	if err := snapshot.Save(s, path); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Secrets["x"] != "val" {
		t.Errorf("expected val, got %q", loaded.Secrets["x"])
	}
}

func TestLoad_InvalidFile_ReturnsError(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRestore_WritesSecretsBack(t *testing.T) {
	src := newStub(t)
	_ = src.Set("k1", "v1")
	s, _ := snapshot.Take(src, []string{"k1"})

	dst := newStub(t)
	if err := snapshot.Restore(dst, s); err != nil {
		t.Fatalf("restore: %v", err)
	}
	got, _ := dst.Get("k1")
	if got != "v1" {
		t.Errorf("expected v1, got %q", got)
	}
}

func TestSave_InvalidPath_ReturnsError(t *testing.T) {
	s := &snapshot.Snapshot{Secrets: map[string]string{}}
	err := snapshot.Save(s, filepath.Join(t.TempDir(), "no", "dir", "snap.json"))
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
	_ = os.Remove("snap.json")
}
