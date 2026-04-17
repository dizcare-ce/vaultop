package validate

import (
	"testing"
)

func TestValidate_Passes(t *testing.T) {
	r := Validate("key", "s3cr3tValue!", Rule{MinLength: 8, MaxLength: 32, Pattern: `^[a-zA-Z0-9!]+$`})
	if !r.Passed {
		t.Fatalf("expected pass, got errors: %v", r.Errors)
	}
}

func TestValidate_BelowMinLength(t *testing.T) {
	r := Validate("key", "abc", Rule{MinLength: 8})
	if r.Passed {
		t.Fatal("expected failure for short value")
	}
	if len(r.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(r.Errors))
	}
}

func TestValidate_ExceedsMaxLength(t *testing.T) {
	r := Validate("key", "toolongvalue", Rule{MaxLength: 5})
	if r.Passed {
		t.Fatal("expected failure for long value")
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	r := Validate("key", "hello world", Rule{Pattern: `^\S+$`})
	if r.Passed {
		t.Fatal("expected failure for pattern mismatch")
	}
}

func TestValidate_InvalidPattern(t *testing.T) {
	r := Validate("key", "value", Rule{Pattern: `[invalid`})
	if r.Passed {
		t.Fatal("expected failure for invalid pattern")
	}
}

func TestValidateAll_AllPass(t *testing.T) {
	secrets := map[string]string{"db_pass": "strongPass1", "api_key": "abcdefgh"}
	rules := map[string]Rule{
		"db_pass": {MinLength: 8},
		"api_key": {MinLength: 8},
	}
	_, err := ValidateAll(secrets, rules)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateAll_SomeFail(t *testing.T) {
	secrets := map[string]string{"db_pass": "weak", "api_key": "abcdefgh"}
	rules := map[string]Rule{
		"db_pass": {MinLength: 8},
		"api_key": {MinLength: 8},
	}
	results, err := ValidateAll(secrets, rules)
	if err == nil {
		t.Fatal("expected error for failed validation")
	}
	var failed int
	for _, r := range results {
		if !r.Passed {
			failed++
		}
	}
	if failed != 1 {
		t.Fatalf("expected 1 failure, got %d", failed)
	}
}
