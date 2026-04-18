package access_test

import (
	"testing"

	"github.com/vaultop/vaultop/internal/access"
)

func TestCheck_ReaderCanRead(t *testing.T) {
	if err := access.Check(access.RoleReader, access.OpRead); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCheck_ReaderCannotDelete(t *testing.T) {
	err := access.Check(access.RoleReader, access.OpDelete)
	if err != access.ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestCheck_WriterCanRotate(t *testing.T) {
	if err := access.Check(access.RoleWriter, access.OpRotate); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCheck_WriterCannotDelete(t *testing.T) {
	err := access.Check(access.RoleWriter, access.OpDelete)
	if err != access.ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestCheck_AdminCanDoAll(t *testing.T) {
	ops := []access.Op{access.OpRead, access.OpWrite, access.OpDelete, access.OpRotate}
	for _, op := range ops {
		if err := access.Check(access.RoleAdmin, op); err != nil {
			t.Fatalf("admin should be allowed %s, got %v", op, err)
		}
	}
}

func TestCheck_UnknownRole(t *testing.T) {
	err := access.Check("ghost", access.OpRead)
	if err != access.ErrUnknownRole {
		t.Fatalf("expected ErrUnknownRole, got %v", err)
	}
}

func TestPermissions_ReturnsCorrectOps(t *testing.T) {
	ops, err := access.Permissions(access.RoleWriter)
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) == 0 {
		t.Fatal("expected non-empty permissions")
	}
}

func TestPermissions_UnknownRole(t *testing.T) {
	_, err := access.Permissions("nobody")
	if err != access.ErrUnknownRole {
		t.Fatalf("expected ErrUnknownRole, got %v", err)
	}
}

func TestIsValid(t *testing.T) {
	if !access.IsValid(access.RoleAdmin) {
		t.Fatal("admin should be valid")
	}
	if access.IsValid("superuser") {
		t.Fatal("superuser should not be valid")
	}
}
