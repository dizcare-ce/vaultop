package search_test

import (
	"context"
	"testing"

	"github.com/your-org/vaultop/internal/search"
)

type stubProvider struct {
	data map[string]string
}

func (s *stubProvider) List(_ context.Context) ([]string, error) {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys, nil
}

func (s *stubProvider) Get(_ context.Context, key string) (string, error) {
	return s.data[key], nil
}

func newStub() *stubProvider {
	return &stubProvider{data: map[string]string{
		"app/db/password": "s3cr3t",
		"app/api/key":     "abc123",
		"infra/tls/cert":  "certdata",
	}}
}

func TestFind_NoFilter_ReturnsAll(t *testing.T) {
	p := newStub()
	res, err := search.Find(context.Background(), p, search.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 3 {
		t.Fatalf("expected 3 results, got %d", len(res))
	}
}

func TestFind_PrefixFilter(t *testing.T) {
	p := newStub()
	res, err := search.Find(context.Background(), p, search.Options{Prefix: "app/"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
}

func TestFind_ContainsFilter(t *testing.T) {
	p := newStub()
	res, err := search.Find(context.Background(), p, search.Options{Contains: "tls"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 || res[0].Key != "infra/tls/cert" {
		t.Fatalf("unexpected results: %v", res)
	}
}

func TestFind_ValueContains(t *testing.T) {
	p := newStub()
	res, err := search.Find(context.Background(), p, search.Options{ValueContains: "abc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 || res[0].Key != "app/api/key" {
		t.Fatalf("unexpected results: %v", res)
	}
}

func TestFind_NoMatch_ReturnsEmpty(t *testing.T) {
	p := newStub()
	res, err := search.Find(context.Background(), p, search.Options{Prefix: "nonexistent/"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("expected 0 results, got %d", len(res))
	}
}
