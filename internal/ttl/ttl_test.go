package ttl_test

import (
	"testing"
	"time"

	"github.com/vaultop/internal/ttl"
)

var epoch = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func TestSet_And_Get(t *testing.T) {
	m := ttl.New()
	m.Set("db/pass", time.Hour, epoch)
	e, ok := m.Get("db/pass")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if !e.ExpiresAt.Equal(epoch.Add(time.Hour)) {
		t.Errorf("unexpected expiry: %v", e.ExpiresAt)
	}
}

func TestGet_Missing(t *testing.T) {
	m := ttl.New()
	_, ok := m.Get("missing")
	if ok {
		t.Fatal("expected no entry")
	}
}

func TestIsExpired_True(t *testing.T) {
	m := ttl.New()
	m.Set("k", time.Minute, epoch)
	now := epoch.Add(2 * time.Minute)
	if err := m.Check("k", now); err != ttl.ErrExpired {
		t.Errorf("expected ErrExpired, got %v", err)
	}
}

func TestIsExpired_False(t *testing.T) {
	m := ttl.New()
	m.Set("k", time.Hour, epoch)
	if err := m.Check("k", epoch); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCheck_MissingKey_NoError(t *testing.T) {
	m := ttl.New()
	if err := m.Check("nope", epoch); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestExpired_ReturnsExpiredKeys(t *testing.T) {
	m := ttl.New()
	m.Set("a", time.Minute, epoch)
	m.Set("b", time.Hour, epoch)
	now := epoch.Add(5 * time.Minute)
	expired := m.Expired(now)
	if len(expired) != 1 || expired[0] != "a" {
		t.Errorf("unexpected expired keys: %v", expired)
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	m := ttl.New()
	m.Set("x", time.Hour, epoch)
	m.Delete("x")
	_, ok := m.Get("x")
	if ok {
		t.Fatal("expected entry to be deleted")
	}
}
