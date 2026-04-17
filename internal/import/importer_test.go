package importer

import (
	"strings"
	"testing"

	"github.com/vaultop/internal/provider"
)

func newStub(t *testing.T) provider.Provider {
	t.Helper()
	p, err := provider.New(provider.Type("stub"), nil)
	if err != nil {
		t.Fatalf("new stub provider: %v", err)
	}
	return p
}

func TestFormat_IsValid(t *testing.T) {
	if !FormatJSON.IsValid() {
		t.Error("json should be valid")
	}
	if !FormatEnv.IsValid() {
		t.Error("env should be valid")
	}
	if Format("xml").IsValid() {
		t.Error("xml should not be valid")
	}
}

func TestImport_JSON(t *testing.T) {
	p := newStub(t)
	r := strings.NewReader(`{"DB_PASS":"secret","API_KEY":"abc123"}`)
	got, err := Import(r, Options{Provider: p, Format: FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_PASS"] != "secret" {
		t.Errorf("expected 'secret', got %q", got["DB_PASS"])
	}
	val, _ := p.Get("API_KEY")
	if val != "abc123" {
		t.Errorf("expected 'abc123' in provider, got %q", val)
	}
}

func TestImport_Env(t *testing.T) {
	p := newStub(t)
	input := "# comment\nDB_PASS=secret\nAPI_KEY=abc123\n"
	got, err := Import(strings.NewReader(input), Options{Provider: p, Format: FormatEnv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 entries, got %d", len(got))
	}
}

func TestImport_DryRun_DoesNotPersist(t *testing.T) {
	p := newStub(t)
	r := strings.NewReader(`{"KEY":"val"}`)
	_, err := Import(r, Options{Provider: p, Format: FormatJSON, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys, _ := p.List()
	if len(keys) != 0 {
		t.Errorf("dry run should not write to provider")
	}
}

func TestImport_UnsupportedFormat(t *testing.T) {
	p := newStub(t)
	_, err := Import(strings.NewReader(""), Options{Provider: p, Format: "toml"})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestImport_InvalidEnvLine(t *testing.T) {
	p := newStub(t)
	_, err := Import(strings.NewReader("BADLINE"), Options{Provider: p, Format: FormatEnv})
	if err == nil {
		t.Error("expected error for invalid env line")
	}
}
