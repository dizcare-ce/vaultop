package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/vaultop/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "vaultop.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return path
}

func TestLoad_ValidConfig(t *testing.T) {
	content := `
version: "1"
providers:
  primary:
    type: aws
    region: us-east-1
defaults:
  rotation_interval: 24h
  dry_run: false
`
	path := writeTempConfig(t, content)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Version != "1" {
		t.Errorf("expected version 1, got %s", cfg.Version)
	}
	if cfg.Providers["primary"].Type != config.ProviderAWS {
		t.Errorf("expected aws provider")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidProvider(t *testing.T) {
	content := `
version: "1"
providers:
  bad:
    type: unknown
`
	path := writeTempConfig(t, content)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for unknown provider type")
	}
}

func TestLoad_MissingVersion(t *testing.T) {
	content := `
providers:
  primary:
    type: gcp
`
	path := writeTempConfig(t, content)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing version")
	}
}
