package env_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seanpk/vaultop/internal/env"
)

// stubProvider is an in-memory provider used by tests.
type stubProvider struct {
	data map[string]string
}

func (s *stubProvider) Get(_ context.Context, key string) (string, error) {
	v, ok := s.data[key]
	if !ok {
		return "", errors.New("not found: " + key)
	}
	return v, nil
}

func newStub(kv ...string) *stubProvider {
	p := &stubProvider{data: make(map[string]string)}
	for i := 0; i+1 < len(kv); i += 2 {
		p.data[kv[i]] = kv[i+1]
	}
	return p
}

func TestExpand_NoPlaceholder(t *testing.T) {
	p := newStub()
	got, err := env.Expand(context.Background(), p, "no secrets here")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "no secrets here" {
		t.Fatalf("expected unchanged string, got %q", got)
	}
}

func TestExpand_SinglePlaceholder(t *testing.T) {
	p := newStub("db/password", "s3cr3t")
	got, err := env.Expand(context.Background(), p, "pass=${secret:db/password}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "pass=s3cr3t" {
		t.Fatalf("got %q", got)
	}
}

func TestExpand_MultiplePlaceholders(t *testing.T) {
	p := newStub("user", "admin", "pass", "hunter2")
	input := "${secret:user}:${secret:pass}"
	got, err := env.Expand(context.Background(), p, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "admin:hunter2" {
		t.Fatalf("got %q", got)
	}
}

func TestExpand_MissingKey_ReturnsError(t *testing.T) {
	p := newStub()
	_, err := env.Expand(context.Background(), p, "${secret:missing}")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestExpandMap_SubstitutesAllValues(t *testing.T) {
	p := newStub("tok", "abc123")
	m := map[string]string{
		"TOKEN": "${secret:tok}",
		"PLAIN": "no-ref",
	}
	out, err := env.ExpandMap(context.Background(), p, m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TOKEN"] != "abc123" {
		t.Fatalf("TOKEN: got %q", out["TOKEN"])
	}
	if out["PLAIN"] != "no-ref" {
		t.Fatalf("PLAIN: got %q", out["PLAIN"])
	}
}

func TestExpandMap_MissingKey_ReturnsError(t *testing.T) {
	p := newStub()
	_, err := env.ExpandMap(context.Background(), p, map[string]string{"X": "${secret:gone}"})
	if err == nil {
		t.Fatal("expected error")
	}
}
