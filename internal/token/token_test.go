package token

import (
	"testing"
	"time"
)

var testSecret = []byte("supersecretkey1234567")

func newManager(t *testing.T, ttl time.Duration) *Manager {
	t.Helper()
	m, err := New(testSecret, ttl)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return m
}

func TestNew_ShortSecret_ReturnsError(t *testing.T) {
	_, err := New([]byte("short"), time.Minute)
	if err == nil {
		t.Fatal("expected error for short secret")
	}
}

func TestNew_ZeroTTL_ReturnsError(t *testing.T) {
	_, err := New(testSecret, 0)
	if err == nil {
		t.Fatal("expected error for zero TTL")
	}
}

func TestIssue_And_Validate_RoundTrip(t *testing.T) {
	m := newManager(t, time.Minute)
	tok, err := m.Issue("svc-a")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	subject, err := m.Validate(tok.Value)
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if subject != "svc-a" {
		t.Errorf("subject = %q; want %q", subject, "svc-a")
	}
}

func TestValidate_ExpiredToken_ReturnsErrExpired(t *testing.T) {
	m := newManager(t, time.Millisecond)
	tok, _ := m.Issue("svc-b")
	time.Sleep(5 * time.Millisecond)
	_, err := m.Validate(tok.Value)
	if err != ErrExpired {
		t.Errorf("err = %v; want ErrExpired", err)
	}
}

func TestValidate_TamperedToken_ReturnsErrInvalid(t *testing.T) {
	m := newManager(t, time.Minute)
	tok, _ := m.Issue("svc-c")
	tampered := tok.Value[:len(tok.Value)-4] + "XXXX"
	_, err := m.Validate(tampered)
	if err != ErrInvalid {
		t.Errorf("err = %v; want ErrInvalid", err)
	}
}

func TestValidate_GarbageInput_ReturnsErrInvalid(t *testing.T) {
	m := newManager(t, time.Minute)
	_, err := m.Validate("not-a-token")
	if err == nil {
		t.Fatal("expected error for garbage input")
	}
}

func TestToken_IsExpired(t *testing.T) {
	past := Token{ExpiresAt: time.Now().Add(-time.Second)}
	if !past.IsExpired() {
		t.Error("expected past token to be expired")
	}
	future := Token{ExpiresAt: time.Now().Add(time.Minute)}
	if future.IsExpired() {
		t.Error("expected future token to not be expired")
	}
}
