package presign_test

import (
	"errors"
	"testing"
	"time"

	"github.com/yourusername/vaultop/internal/presign"
)

func newManager(t *testing.T) *presign.Manager {
	t.Helper()
	m, err := presign.New(presign.Config{
		Secret:     "supersecretkey1234567890abcdef!!",
		DefaultTTL: 5 * time.Minute,
	})
	if err != nil {
		t.Fatalf("presign.New: %v", err)
	}
	return m
}

func TestNew_ShortSecret_ReturnsError(t *testing.T) {
	_, err := presign.New(presign.Config{
		Secret:     "short",
		DefaultTTL: time.Minute,
	})
	if err == nil {
		t.Fatal("expected error for short secret, got nil")
	}
}

func TestNew_ZeroTTL_ReturnsError(t *testing.T) {
	_, err := presign.New(presign.Config{
		Secret:     "supersecretkey1234567890abcdef!!",
		DefaultTTL: 0,
	})
	if err == nil {
		t.Fatal("expected error for zero TTL, got nil")
	}
}

func TestSign_And_Verify_RoundTrip(t *testing.T) {
	m := newManager(t)

	token, err := m.Sign("db/password", nil)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := m.Verify(token)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if claims.Key != "db/password" {
		t.Errorf("claims.Key = %q, want %q", claims.Key, "db/password")
	}
}

func TestSign_WithCustomTTL(t *testing.T) {
	m := newManager(t)

	token, err := m.Sign("api/key", &presign.SignOptions{TTL: 10 * time.Second})
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	claims, err := m.Verify(token)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if claims.Key != "api/key" {
		t.Errorf("claims.Key = %q, want %q", claims.Key, "api/key")
	}
}

func TestVerify_ExpiredToken_ReturnsError(t *testing.T) {
	m, err := presign.New(presign.Config{
		Secret:     "supersecretkey1234567890abcdef!!",
		DefaultTTL: time.Millisecond,
	})
	if err != nil {
		t.Fatalf("presign.New: %v", err)
	}

	token, err := m.Sign("db/password", nil)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	time.Sleep(5 * time.Millisecond)

	_, err = m.Verify(token)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
	if !errors.Is(err, presign.ErrExpired) {
		t.Errorf("expected ErrExpired, got %v", err)
	}
}

func TestVerify_TamperedToken_ReturnsError(t *testing.T) {
	m := newManager(t)

	token, err := m.Sign("db/password", nil)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	tampered := token[:len(token)-4] + "XXXX"
	_, err = m.Verify(tampered)
	if err == nil {
		t.Fatal("expected error for tampered token, got nil")
	}
}

func TestVerify_EmptyToken_ReturnsError(t *testing.T) {
	m := newManager(t)

	_, err := m.Verify("")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}

func TestSign_WithMetadata(t *testing.T) {
	m := newManager(t)

	opts := &presign.SignOptions{
		TTL:      time.Minute,
		Metadata: map[string]string{"actor": "ci-bot", "env": "staging"},
	}

	token, err := m.Sign("infra/token", opts)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	claims, err := m.Verify(token)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if claims.Metadata["actor"] != "ci-bot" {
		t.Errorf("metadata actor = %q, want %q", claims.Metadata["actor"], "ci-bot")
	}
	if claims.Metadata["env"] != "staging" {
		t.Errorf("metadata env = %q, want %q", claims.Metadata["env"], "staging")
	}
}
