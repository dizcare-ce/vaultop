package access_test

import (
	"testing"

	"github.com/vaultop/vaultop/internal/access"
)

func basePolicy() *access.Policy {
	return &access.Policy{
		Entries: []access.Entry{
			{Identity: "alice", Role: access.RoleAdmin},
			{Identity: "bob", Role: access.RoleReader},
			{Identity: "ci", Role: access.RoleWriter},
		},
	}
}

func TestPolicy_Lookup_Found(t *testing.T) {
	p := basePolicy()
	role, err := p.Lookup("alice")
	if err != nil {
		t.Fatal(err)
	}
	if role != access.RoleAdmin {
		t.Fatalf("expected admin, got %s", role)
	}
}

func TestPolicy_Lookup_NotFound(t *testing.T) {
	p := basePolicy()
	_, err := p.Lookup("unknown")
	if err == nil {
		t.Fatal("expected error for unknown identity")
	}
}

func TestPolicy_Allow_Permitted(t *testing.T) {
	p := basePolicy()
	if err := p.Allow("ci", access.OpWrite); err != nil {
		t.Fatalf("expected allowed, got %v", err)
	}
}

func TestPolicy_Allow_Denied(t *testing.T) {
	p := basePolicy()
	if err := p.Allow("bob", access.OpDelete); err == nil {
		t.Fatal("expected denial for reader attempting delete")
	}
}

func TestPolicy_Validate_Valid(t *testing.T) {
	if err := basePolicy().Validate(); err != nil {
		t.Fatalf("expected valid policy, got %v", err)
	}
}

func TestPolicy_Validate_UnknownRole(t *testing.T) {
	p := &access.Policy{
		Entries: []access.Entry{{Identity: "x", Role: "superadmin"}},
	}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for unknown role")
	}
}

func TestPolicy_Validate_EmptyIdentity(t *testing.T) {
	p := &access.Policy{
		Entries: []access.Entry{{Identity: "", Role: access.RoleReader}},
	}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for empty identity")
	}
}
