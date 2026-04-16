package provider

import (
	"testing"
)

func TestTypeIsValid(t *testing.T) {
	valid := []Type{AWS, GCP, Azure, Vault}
	for _, p := range valid {
		if !p.IsValid() {
			t.Errorf("expected %q to be valid", p)
		}
	}
	if Type("unknown").IsValid() {
		t.Error("expected 'unknown' to be invalid")
	}
}

func TestNew_UnsupportedProvider(t *testing.T) {
	_, err := New("bogus", nil)
	if err == nil {
		t.Fatal("expected error for unsupported provider")
	}
}

func TestStubProvider_SetGet(t *testing.T) {
	p, err := New(AWS, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := p.SetSecret("db/password", "s3cr3t"); err != nil {
		t.Fatalf("SetSecret: %v", err)
	}
	v, err := p.GetSecret("db/password")
	if err != nil {
		t.Fatalf("GetSecret: %v", err)
	}
	if v != "s3cr3t" {
		t.Errorf("got %q, want %q", v, "s3cr3t")
	}
}

func TestStubProvider_DeleteAndList(t *testing.T) {
	p, _ := New(Vault, nil)
	_ = p.SetSecret("app/key1", "v1")
	_ = p.SetSecret("app/key2", "v2")
	_ = p.SetSecret("other", "v3")

	list, err := p.ListSecrets("app/")
	if err != nil {
		t.Fatalf("ListSecrets: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(list))
	}

	if err := p.DeleteSecret("app/key1"); err != nil {
		t.Fatalf("DeleteSecret: %v", err)
	}
	if _, err := p.GetSecret("app/key1"); err == nil {
		t.Error("expected error after deletion")
	}
}
