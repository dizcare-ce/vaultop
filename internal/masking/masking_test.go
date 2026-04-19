package masking_test

import (
	"testing"

	"github.com/vaultop/vaultop/internal/masking"
)

func TestApply_Full(t *testing.T) {
	cfg := masking.DefaultConfig()
	if got := masking.Apply("supersecret", cfg); got != "***" {
		t.Fatalf("expected *** got %s", got)
	}
}

func TestApply_Full_EmptyValue(t *testing.T) {
	cfg := masking.DefaultConfig()
	if got := masking.Apply("", cfg); got != "***" {
		t.Fatalf("expected *** got %s", got)
	}
}

func TestApply_Partial_ShowsSuffix(t *testing.T) {
	cfg := masking.Config{Mode: masking.ModePartial, Mask: "***", VisibleLen: 4}
	got := masking.Apply("abcdefgh", cfg)
	if got != "***efgh" {
		t.Fatalf("expected ***efgh got %s", got)
	}
}

func TestApply_Partial_ShortValue(t *testing.T) {
	cfg := masking.Config{Mode: masking.ModePartial, Mask: "***", VisibleLen: 6}
	got := masking.Apply("abc", cfg)
	if got != "***" {
		t.Fatalf("expected *** got %s", got)
	}
}

func TestApply_Prefix_ShowsPrefix(t *testing.T) {
	cfg := masking.Config{Mode: masking.ModePrefix, Mask: "***", VisibleLen: 3}
	got := masking.Apply("abcdefgh", cfg)
	if got != "abc***" {
		t.Fatalf("expected abc*** got %s", got)
	}
}

func TestApplyMap_MasksAllValues(t *testing.T) {
	cfg := masking.DefaultConfig()
	secrets := map[string]string{"db_pass": "hunter2", "api_key": "xyz123"}
	out := masking.ApplyMap(secrets, cfg)
	for k, v := range out {
		if v != "***" {
			t.Errorf("key %s: expected *** got %s", k, v)
		}
	}
}

func TestContainsSensitive_True(t *testing.T) {
	if !masking.ContainsSensitive("error: bad password hunter2 rejected", []string{"hunter2"}) {
		t.Fatal("expected true")
	}
}

func TestContainsSensitive_False(t *testing.T) {
	if masking.ContainsSensitive("everything is fine", []string{"hunter2"}) {
		t.Fatal("expected false")
	}
}

func TestContainsSensitive_EmptySecret_Skipped(t *testing.T) {
	if masking.ContainsSensitive("some log line", []string{""}) {
		t.Fatal("empty secret should not match")
	}
}
