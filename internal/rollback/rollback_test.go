package rollback_test

import (
	"context"
	"errors"
	"testing"

	"github.com/vaultop/internal/audit"
	"github.com/vaultop/internal/history"
	"github.com/vaultop/internal/rollback"
	"github.com/vaultop/internal/snapshot"
)

// minimal in-memory stub provider
type stubProvider struct {
	data    map[string]string
	setFail bool
}

func newStub(data map[string]string) *stubProvider {
	return &stubProvider{data: data}
}

func (s *stubProvider) Get(_ context.Context, key string) (string, error) {
	v, ok := s.data[key]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

func (s *stubProvider) Set(_ context.Context, key, val string) error {
	if s.setFail {
		return errors.New("set error")
	}
	s.data[key] = val
	return nil
}

func (s *stubProvider) Delete(_ context.Context, key string) error { delete(s.data, key); return nil }
func (s *stubProvider) List(_ context.Context) ([]string, error)   { return nil, nil }

func buildOpts(t *testing.T, p *stubProvider) rollback.Options {
	t.Helper()
	h, err := history.Load(t.TempDir() + "/hist.json")
	if err != nil {
		t.Fatal(err)
	}
	return rollback.Options{
		Provider: p,
		History:  h,
		Audit:    audit.Noop{},
	}
}

func TestRun_RestoresValues(t *testing.T) {
	p := newStub(map[string]string{"db/pass": "old"})
	snap := snapshot.Snapshot{"db/pass": "restored"}
	opts := buildOpts(t, p)

	results := rollback.Run(context.Background(), snap, opts)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if p.data["db/pass"] != "restored" {
		t.Errorf("expected restored, got %s", p.data["db/pass"])
	}
}

func TestRun_DryRun_DoesNotPersist(t *testing.T) {
	p := newStub(map[string]string{"key": "original"})
	snap := snapshot.Snapshot{"key": "new-value"}
	opts := buildOpts(t, p)
	opts.DryRun = true

	rollback.Run(context.Background(), snap, opts)

	if p.data["key"] != "original" {
		t.Errorf("dry run should not modify provider")
	}
}

func TestRun_SetError_RecordsFailure(t *testing.T) {
	p := newStub(map[string]string{})
	p.setFail = true
	snap := snapshot.Snapshot{"key": "val"}
	opts := buildOpts(t, p)

	results := rollback.Run(context.Background(), snap, opts)

	if !rollback.AnyFailed(results) {
		t.Error("expected failure to be recorded")
	}
}

func TestAnyFailed_False(t *testing.T) {
	results := []rollback.Result{{Key: "k", Err: nil}}
	if rollback.AnyFailed(results) {
		t.Error("expected no failures")
	}
}
