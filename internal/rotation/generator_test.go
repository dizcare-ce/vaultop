package rotation_test

import (
	"strings"
	"testing"

	"github.com/vaultop/internal/rotation"
)

func TestDefaultGenerator_ReturnsNonEmptyString(t *testing.T) {
	val, err := rotation.DefaultGenerator("any-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val == "" {
		t.Error("expected non-empty value")
	}
}

func TestDefaultGenerator_UniqueValues(t *testing.T) {
	a, _ := rotation.DefaultGenerator("id")
	b, _ := rotation.DefaultGenerator("id")
	if a == b {
		t.Error("expected unique values across calls")
	}
}

func TestFixedGenerator_ReturnsFixed(t *testing.T) {
	gen := rotation.FixedGenerator("static-secret")
	for _, id := range []string{"a", "b", "c"} {
		val, err := gen(id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "static-secret" {
			t.Errorf("expected 'static-secret', got %q", val)
		}
	}
}

func TestErrorGenerator_ReturnsError(t *testing.T) {
	gen := rotation.ErrorGenerator("something went wrong")
	_, err := gen("id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "something went wrong") {
		t.Errorf("unexpected error message: %v", err)
	}
}
