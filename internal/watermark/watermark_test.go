package watermark_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultop/internal/watermark"
)

var validSecret = []byte("super-secret-key-32bytes-padding!")

func newManager(t *testing.T) *watermark.Manager {
	t.Helper()
	m, err := watermark.New(validSecret)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return m
}

func TestNew_ShortSecret_ReturnsError(t *testing.T) {
	_, err := watermark.New([]byte("short"))
	if err == nil {
		t.Fatal("expected error for short secret, got nil")
	}
}

func TestApply_And_Verify_RoundTrip(t *testing.T) {
	m := newManager(t)
	marked := m.Apply("db/password", "s3cr3t")
	got, err := m.Verify("db/password", marked)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if got != "s3cr3t" {
		t.Fatalf("expected %q, got %q", "s3cr3t", got)
	}
}

func TestVerify_TamperedValue_ReturnsError(t *testing.T) {
	m := newManager(t)
	marked := m.Apply("db/password", "original")
	tampered := strings.Replace(marked, "original", "hacked", 1)
	_, err := m.Verify("db/password", tampered)
	if err != watermark.ErrInvalidWatermark {
		t.Fatalf("expected ErrInvalidWatermark, got %v", err)
	}
}

func TestVerify_WrongKey_ReturnsError(t *testing.T) {
	m := newManager(t)
	marked := m.Apply("key/a", "value")
	_, err := m.Verify("key/b", marked)
	if err != watermark.ErrInvalidWatermark {
		t.Fatalf("expected ErrInvalidWatermark, got %v", err)
	}
}

func TestVerify_NoSeparator_ReturnsError(t *testing.T) {
	m := newManager(t)
	_, err := m.Verify("key", "plainvalue")
	if err != watermark.ErrInvalidWatermark {
		t.Fatalf("expected ErrInvalidWatermark, got %v", err)
	}
}

func TestIsMarked_True(t *testing.T) {
	m := newManager(t)
	marked := m.Apply("k", "v")
	if !watermark.IsMarked(marked) {
		t.Fatal("expected IsMarked to return true")
	}
}

func TestIsMarked_False(t *testing.T) {
	if watermark.IsMarked("plainvalue") {
		t.Fatal("expected IsMarked to return false for plain value")
	}
}
