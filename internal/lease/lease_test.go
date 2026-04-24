package lease

import (
	"testing"
	"time"
)

// fixedClock returns a clock function pinned to t.
func fixedClock(t time.Time) func() time.Time { return func() time.Time { return t } }

func newManager(now time.Time) *Manager {
	m := New()
	m.now = fixedClock(now)
	return m
}

func TestGrant_CreatesLease(t *testing.T) {
	now := time.Now()
	m := newManager(now)
	l := m.Grant("db/password", time.Hour)
	if l.Key != "db/password" {
		t.Fatalf("expected key db/password, got %s", l.Key)
	}
	if !l.IssuedAt.Equal(now) {
		t.Fatalf("unexpected issued time")
	}
	if l.IsExpired(now) {
		t.Fatal("newly granted lease should not be expired")
	}
}

func TestGet_ExpiredLease_ReturnsErrExpired(t *testing.T) {
	now := time.Now()
	m := newManager(now)
	m.Grant("k", time.Millisecond)
	m.now = fixedClock(now.Add(time.Second))
	_, err := m.Get("k")
	if err != ErrExpired {
		t.Fatalf("expected ErrExpired, got %v", err)
	}
}

func TestGet_MissingLease_ReturnsErrNotFound(t *testing.T) {
	m := newManager(time.Now())
	_, err := m.Get("missing")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRenew_ExtendsExpiry(t *testing.T) {
	now := time.Now()
	m := newManager(now)
	m.Grant("k", time.Minute)
	l, err := m.Renew("k", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Renewals != 1 {
		t.Fatalf("expected 1 renewal, got %d", l.Renewals)
	}
	if l.IsExpired(now.Add(90 * time.Minute)) {
		t.Fatal("lease should still be valid 90 minutes after renewal")
	}
}

func TestRenew_UnknownKey_ReturnsErrNotFound(t *testing.T) {
	m := newManager(time.Now())
	_, err := m.Renew("ghost", time.Hour)
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRevoke_RemovesLease(t *testing.T) {
	m := newManager(time.Now())
	m.Grant("k", time.Hour)
	m.Revoke("k")
	_, err := m.Get("k")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound after revoke, got %v", err)
	}
}

func TestActive_ExcludesExpired(t *testing.T) {
	now := time.Now()
	m := newManager(now)
	m.Grant("live", time.Hour)
	m.Grant("dead", time.Millisecond)
	m.now = fixedClock(now.Add(time.Second))
	active := m.Active()
	if len(active) != 1 || active[0].Key != "live" {
		t.Fatalf("expected only 'live' lease, got %v", active)
	}
}

func TestPurge_RemovesExpiredLeases(t *testing.T) {
	now := time.Now()
	m := newManager(now)
	m.Grant("a", time.Millisecond)
	m.Grant("b", time.Millisecond)
	m.Grant("c", time.Hour)
	m.now = fixedClock(now.Add(time.Second))
	n := m.Purge()
	if n != 2 {
		t.Fatalf("expected 2 purged, got %d", n)
	}
	if len(m.Active()) != 1 {
		t.Fatal("expected 1 active lease after purge")
	}
}
