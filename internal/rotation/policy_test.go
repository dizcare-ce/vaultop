package rotation

import (
	"testing"
	"time"
)

func TestPolicy_Validate_Valid(t *testing.T) {
	p := Policy{Key: "db/password", Interval: 24 * time.Hour, Length: 32}
	if err := p.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPolicy_Validate_EmptyKey(t *testing.T) {
	p := Policy{Key: "", Interval: 24 * time.Hour}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestPolicy_Validate_NegativeLength(t *testing.T) {
	p := Policy{Key: "api/token", Length: -1}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for negative length")
	}
}

func TestPolicy_EffectiveLength_UsesDefault(t *testing.T) {
	p := Policy{Key: "k", Length: 0}
	if got := p.EffectiveLength(); got != DefaultSecretLength {
		t.Fatalf("expected %d, got %d", DefaultSecretLength, got)
	}
}

func TestPolicy_EffectiveLength_UsesExplicit(t *testing.T) {
	p := Policy{Key: "k", Length: 64}
	if got := p.EffectiveLength(); got != 64 {
		t.Fatalf("expected 64, got %d", got)
	}
}
