package fingerprint_test

import (
	"strings"
	"testing"

	"github.com/vaultop/internal/fingerprint"
)

func TestAlgorithm_IsValid(t *testing.T) {
	if !fingerprint.AlgorithmSHA256.IsValid() {
		t.Error("expected sha256 to be valid")
	}
	if !fingerprint.AlgorithmPrefix.IsValid() {
		t.Error("expected prefix to be valid")
	}
	if fingerprint.Algorithm("md5").IsValid() {
		t.Error("expected md5 to be invalid")
	}
}

func TestOf_SHA256_NonEmpty(t *testing.T) {
	fp, err := fingerprint.Of("supersecret", fingerprint.AlgorithmSHA256)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(fp, "sha256:") {
		t.Errorf("expected sha256 prefix, got %q", fp)
	}
}

func TestOf_SHA256_Deterministic(t *testing.T) {
	a, _ := fingerprint.Of("value", fingerprint.AlgorithmSHA256)
	b, _ := fingerprint.Of("value", fingerprint.AlgorithmSHA256)
	if a != b {
		t.Errorf("expected identical fingerprints, got %q and %q", a, b)
	}
}

func TestOf_SHA256_DifferentValues(t *testing.T) {
	a, _ := fingerprint.Of("abc", fingerprint.AlgorithmSHA256)
	b, _ := fingerprint.Of("xyz", fingerprint.AlgorithmSHA256)
	if a == b {
		t.Error("expected different fingerprints for different values")
	}
}

func TestOf_EmptyValue_ReturnsError(t *testing.T) {
	_, err := fingerprint.Of("", fingerprint.AlgorithmSHA256)
	if err == nil {
		t.Fatal("expected error for empty value")
	}
	if err != fingerprint.ErrEmptyValue {
		t.Errorf("expected ErrEmptyValue, got %v", err)
	}
}

func TestOf_UnsupportedAlgorithm_ReturnsError(t *testing.T) {
	_, err := fingerprint.Of("val", fingerprint.Algorithm("unknown"))
	if err == nil {
		t.Fatal("expected error for unsupported algorithm")
	}
}

func TestOf_Prefix_ContainsLength(t *testing.T) {
	fp, err := fingerprint.Of("hello world", fingerprint.AlgorithmPrefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fp, "len=11") {
		t.Errorf("expected length in prefix fingerprint, got %q", fp)
	}
}

func TestMapOf_ReturnsAllKeys(t *testing.T) {
	secrets := map[string]string{
		"db_pass": "hunter2",
		"api_key": "abc123",
	}
	out, err := fingerprint.MapOf(secrets, fingerprint.AlgorithmSHA256)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k := range secrets {
		if _, ok := out[k]; !ok {
			t.Errorf("missing key %q in output", k)
		}
	}
}

func TestMapOf_EmptyValue_ReturnsError(t *testing.T) {
	secrets := map[string]string{"key": ""}
	_, err := fingerprint.MapOf(secrets, fingerprint.AlgorithmSHA256)
	if err == nil {
		t.Fatal("expected error for empty secret value")
	}
}

func TestEqual_CaseInsensitive(t *testing.T) {
	if !fingerprint.Equal("SHA256:ABCD", "sha256:abcd") {
		t.Error("expected Equal to be case-insensitive")
	}
}

func TestEqual_Different(t *testing.T) {
	if fingerprint.Equal("sha256:aaa", "sha256:bbb") {
		t.Error("expected unequal fingerprints to return false")
	}
}
