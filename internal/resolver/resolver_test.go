package resolver_test

import (
	"errors"
	"testing"

	"github.com/yourorg/vaultop/internal/resolver"
)

// stubProvider is an in-memory provider for tests.
type stubProvider struct {
	data map[string]string
}

func (s *stubProvider) Get(key string) (string, error) {
	v, ok := s.data[key]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

func newStub(data map[string]string) *stubProvider {
	return &stubProvider{data: data}
}

func TestResolve_DirectKey(t *testing.T) {
	p := newStub(map[string]string{"db/password": "s3cr3t"})
	r := resolver.New(p, resolver.Config{})

	val, err := r.Resolve("db/password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "s3cr3t" {
		t.Errorf("expected s3cr3t, got %s", val)
	}
}

func TestResolve_Alias(t *testing.T) {
	p := newStub(map[string]string{"prod/db/password": "prod-secret"})
	r := resolver.New(p, resolver.Config{
		Aliases: map[string]string{"db_pass": "prod/db/password"},
	})

	val, err := r.Resolve("db_pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "prod-secret" {
		t.Errorf("expected prod-secret, got %s", val)
	}
}

func TestResolve_Fallback_UsedWhenPrimaryMissing(t *testing.T) {
	p := newStub(map[string]string{"staging/db/password": "stage-secret"})
	r := resolver.New(p, resolver.Config{
		Fallbacks: map[string][]string{
			"db_pass": {"staging/db/password"},
		},
	})

	val, err := r.Resolve("db_pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "stage-secret" {
		t.Errorf("expected stage-secret, got %s", val)
	}
}

func TestResolve_NotResolved_ReturnsError(t *testing.T) {
	p := newStub(map[string]string{})
	r := resolver.New(p, resolver.Config{})

	_, err := r.Resolve("missing/key")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, resolver.ErrNotResolved) {
		t.Errorf("expected ErrNotResolved, got %v", err)
	}
}

func TestResolve_Alias_WithFallback(t *testing.T) {
	p := newStub(map[string]string{"fallback/key": "fallback-val"})
	r := resolver.New(p, resolver.Config{
		Aliases:   map[string]string{"token": "canonical/token"},
		Fallbacks: map[string][]string{"token": {"fallback/key"}},
	})

	val, err := r.Resolve("token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "fallback-val" {
		t.Errorf("expected fallback-val, got %s", val)
	}
}
