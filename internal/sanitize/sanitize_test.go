package sanitize_test

import (
	"testing"

	"github.com/yourusername/vaultop/internal/sanitize"
)

func TestApply_TrimSpace(t *testing.T) {
	opts := sanitize.DefaultOptions()
	got, err := sanitize.Apply("  hello  ", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello" {
		t.Fatalf("expected 'hello', got %q", got)
	}
}

func TestApply_StripControl(t *testing.T) {
	opts := sanitize.DefaultOptions()
	got, err := sanitize.Apply("hel\x00lo", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello" {
		t.Fatalf("expected 'hello', got %q", got)
	}
}

func TestApply_MaxLength_Truncates(t *testing.T) {
	opts := sanitize.Options{TrimSpace: false, StripControl: false, MaxLength: 4, RejectEmpty: false}
	got, err := sanitize.Apply("abcdefgh", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abcd" {
		t.Fatalf("expected 'abcd', got %q", got)
	}
}

func TestApply_RejectEmpty_ReturnsError(t *testing.T) {
	opts := sanitize.DefaultOptions()
	_, err := sanitize.Apply("   ", opts)
	if err != sanitize.ErrEmptyValue {
		t.Fatalf("expected ErrEmptyValue, got %v", err)
	}
}

func TestApply_RejectEmpty_Disabled(t *testing.T) {
	opts := sanitize.Options{TrimSpace: true, RejectEmpty: false}
	got, err := sanitize.Apply("   ", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestApplyMap_SkipsInvalid(t *testing.T) {
	opts := sanitize.DefaultOptions()
	m := map[string]string{
		"good": "  value  ",
		"bad":  "   ",
	}
	out, errs := sanitize.ApplyMap(m, opts)
	if _, ok := out["good"]; !ok {
		t.Fatal("expected 'good' key in output")
	}
	if out["good"] != "value" {
		t.Fatalf("expected 'value', got %q", out["good"])
	}
	if _, ok := errs["bad"]; !ok {
		t.Fatal("expected error for 'bad' key")
	}
	if _, ok := out["bad"]; ok {
		t.Fatal("'bad' key should not appear in output")
	}
}

func TestApplyMap_AllValid(t *testing.T) {
	opts := sanitize.DefaultOptions()
	m := map[string]string{"a": "foo", "b": "bar"}
	out, errs := sanitize.ApplyMap(m, opts)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}
