package tag_test

import (
	"context"
	"testing"

	"github.com/user/vaultop/internal/tag"
)

// minimal in-memory stub
type stubProvider struct{ data map[string]string }

func newStub() *stubProvider { return &stubProvider{data: map[string]string{}} }
func (s *stubProvider) Get(_ context.Context, key string) (string, error) {
	v, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("not found")
	}
	return v, nil
}
func (s *stubProvider) Set(_ context.Context, key, value string) error {
	s.data[key] = value
	return nil
}
func (s *stubProvider) List(_ context.Context) ([]string, error) {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys, nil
}

import "fmt"

func TestSet_And_Get(t *testing.T) {
	ctx := context.Background()
	p := newStub()

	if err := tag.Set(ctx, p, "db/password", "owner", "team-a"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	v, ok, err := tag.Get(ctx, p, "db/password", "owner")
	if err != nil || !ok || v != "team-a" {
		t.Fatalf("Get: got %q %v %v", v, ok, err)
	}
}

func TestDelete_RemovesTag(t *testing.T) {
	ctx := context.Background()
	p := newStub()
	_ = tag.Set(ctx, p, "key", "env", "prod")
	_ = tag.Delete(ctx, p, "key", "env")
	_, ok, _ := tag.Get(ctx, p, "key", "env")
	if ok {
		t.Fatal("expected tag to be deleted")
	}
}

func TestGetAll_Empty(t *testing.T) {
	ctx := context.Background()
	p := newStub()
	tags, err := tag.GetAll(ctx, p, "missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 0 {
		t.Fatalf("expected empty map, got %v", tags)
	}
}

func TestListTagged(t *testing.T) {
	ctx := context.Background()
	p := newStub()
	_ = tag.Set(ctx, p, "secret/a", "team", "x")
	_ = tag.Set(ctx, p, "secret/b", "team", "y")
	keys, err := tag.ListTagged(ctx, p)
	if err != nil {
		t.Fatalf("ListTagged: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("expected 2 tagged secrets, got %d: %v", len(keys), keys)
	}
}
