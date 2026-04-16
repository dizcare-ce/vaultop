package export_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/vaultop/internal/export"
)

func TestFormat_IsValid(t *testing.T) {
	valid := []export.Format{export.FormatJSON, export.FormatEnv, export.FormatDotEnv}
	for _, f := range valid {
		if !f.IsValid() {
			t.Errorf("expected %q to be valid", f)
		}
	}
	if export.Format("xml").IsValid() {
		t.Error("expected 'xml' to be invalid")
	}
}

func TestWrite_JSON(t *testing.T) {
	secrets := map[string]string{"KEY": "value", "OTHER": "data"}
	var buf bytes.Buffer
	if err := export.Write(&buf, secrets, export.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got map[string]string
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	for k, v := range secrets {
		if got[k] != v {
			t.Errorf("key %q: want %q got %q", k, v, got[k])
		}
	}
}

func TestWrite_Env(t *testing.T) {
	secrets := map[string]string{"ALPHA": "one", "BETA": "two"}
	var buf bytes.Buffer
	if err := export.Write(&buf, secrets, export.FormatEnv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ALPHA=one") {
		t.Errorf("missing ALPHA=one in output: %q", out)
	}
	if !strings.Contains(out, "BETA=two") {
		t.Errorf("missing BETA=two in output: %q", out)
	}
}

func TestWrite_DotEnv_AliasForEnv(t *testing.T) {
	secrets := map[string]string{"X": "y"}
	var b1, b2 bytes.Buffer
	_ = export.Write(&b1, secrets, export.FormatEnv)
	_ = export.Write(&b2, secrets, export.FormatDotEnv)
	if b1.String() != b2.String() {
		t.Error("env and dotenv formats should produce identical output")
	}
}

func TestWrite_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := export.Write(&buf, map[string]string{}, export.Format("toml"))
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
