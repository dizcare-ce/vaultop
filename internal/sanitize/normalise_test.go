package sanitize_test

import (
	"testing"

	"github.com/yourusername/vaultop/internal/sanitize"
)

func TestNormaliseKey_Lowercase(t *testing.T) {
	got := sanitize.NormaliseKey("MySecret")
	if got != "mysecret" {
		t.Fatalf("expected 'mysecret', got %q", got)
	}
}

func TestNormaliseKey_ReplacesSpecialChars(t *testing.T) {
	got := sanitize.NormaliseKey("my-secret/key")
	if got != "my_secret_key" {
		t.Fatalf("expected 'my_secret_key', got %q", got)
	}
}

func TestNormaliseKey_CollapseRuns(t *testing.T) {
	got := sanitize.NormaliseKey("my---secret")
	if got != "my_secret" {
		t.Fatalf("expected 'my_secret', got %q", got)
	}
}

func TestNormaliseKey_TrimLeadingTrailing(t *testing.T) {
	got := sanitize.NormaliseKey("-my-secret-")
	if got != "my_secret" {
		t.Fatalf("expected 'my_secret', got %q", got)
	}
}

func TestNormaliseKey_AlreadyNormal(t *testing.T) {
	got := sanitize.NormaliseKey("db_password")
	if got != "db_password" {
		t.Fatalf("expected 'db_password', got %q", got)
	}
}

func TestNormaliseKey_TrimSpace(t *testing.T) {
	got := sanitize.NormaliseKey("  api_key  ")
	if got != "api_key" {
		t.Fatalf("expected 'api_key', got %q", got)
	}
}

func TestNormaliseMap_NormalisesKeys(t *testing.T) {
	m := map[string]string{
		"MY-KEY":   "val1",
		"other/key": "val2",
	}
	out := sanitize.NormaliseMap(m)
	if out["my_key"] != "val1" {
		t.Fatalf("expected val1 for my_key, got %q", out["my_key"])
	}
	if out["other_key"] != "val2" {
		t.Fatalf("expected val2 for other_key, got %q", out["other_key"])
	}
}

func TestNormaliseMap_PreservesValues(t *testing.T) {
	m := map[string]string{"plain": "secret-value"}
	out := sanitize.NormaliseMap(m)
	if out["plain"] != "secret-value" {
		t.Fatalf("value should be unchanged, got %q", out["plain"])
	}
}
