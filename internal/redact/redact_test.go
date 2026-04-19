package redact_test

import (
	"testing"

	"github.com/vaultop/internal/redact"
)

func TestRedact_Full(t *testing.T) {
	got := redact.Redact("supersecret", redact.Options{Mode: redact.ModeFull})
	if got != redact.DefaultMask {
		t.Fatalf("expected %q, got %q", redact.DefaultMask, got)
	}
}

func TestRedact_Full_CustomMask(t *testing.T) {
	got := redact.Redact("supersecret", redact.Options{Mode: redact.ModeFull, Mask: "[REDACTED]"})
	if got != "[REDACTED]" {
		t.Fatalf("unexpected: %q", got)
	}
}

func TestRedact_Partial_ShowsSuffix(t *testing.T) {
	got := redact.Redact("supersecret", redact.Options{Mode: redact.ModePartial, ShowSuffix: 4})
	if got != "***cret" {
		t.Fatalf("expected %q, got %q", "***cret", got)
	}
}

func TestRedact_Partial_ShortValue(t *testing.T) {
	got := redact.Redact("ab", redact.Options{Mode: redact.ModePartial, ShowSuffix: 4})
	if got != redact.DefaultMask {
		t.Fatalf("expected mask for short value, got %q", got)
	}
}

func TestRedact_Hash_ReplacesWithStars(t *testing.T) {
	got := redact.Redact("hello", redact.Options{Mode: redact.ModeHash})
	if got != "*****" {
		t.Fatalf("expected %q, got %q", "*****", got)
	}
}

func TestRedactMap_RedactsAllValues(t *testing.T) {
	m := map[string]string{"key1": "val1", "key2": "val2"}
	out := redact.RedactMap(m, redact.Options{Mode: redact.ModeFull})
	for k, v := range out {
		if v != redact.DefaultMask {
			t.Errorf("key %q: expected mask, got %q", k, v)
		}
	}
	if len(out) != len(m) {
		t.Errorf("map length mismatch")
	}
}

func TestRedactMap_OriginalUnchanged(t *testing.T) {
	m := map[string]string{"k": "secret"}
	redact.RedactMap(m, redact.Options{Mode: redact.ModeFull})
	if m["k"] != "secret" {
		t.Fatal("original map was mutated")
	}
}
