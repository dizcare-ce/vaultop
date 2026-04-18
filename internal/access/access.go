// Package access provides role-based access control for secret operations.
package access

import (
	"errors"
	"slices"
)

// Role represents a named permission level.
type Role string

const (
	RoleReader Role = "reader"
	RoleWriter Role = "writer"
	RoleAdmin  Role = "admin"
)

// Op represents an operation that can be performed on a secret.
type Op string

const (
	OpRead   Op = "read"
	OpWrite  Op = "write"
	OpDelete Op = "delete"
	OpRotate Op = "rotate"
)

var rolePermissions = map[Role][]Op{
	RoleReader: {OpRead},
	RoleWriter: {OpRead, OpWrite, OpRotate},
	RoleAdmin:  {OpRead, OpWrite, OpDelete, OpRotate},
}

// ErrUnauthorized is returned when a role lacks permission for an operation.
var ErrUnauthorized = errors.New("access: unauthorized")

// ErrUnknownRole is returned when the role is not recognised.
var ErrUnknownRole = errors.New("access: unknown role")

// Check returns nil when role is permitted to perform op, otherwise an error.
func Check(role Role, op Op) error {
	perms, ok := rolePermissions[role]
	if !ok {
		return ErrUnknownRole
	}
	if !slices.Contains(perms, op) {
		return ErrUnauthorized
	}
	return nil
}

// Permissions returns the list of operations allowed for role.
func Permissions(role Role) ([]Op, error) {
	perms, ok := rolePermissions[role]
	if !ok {
		return nil, ErrUnknownRole
	}
	out := make([]Op, len(perms))
	copy(out, perms)
	return out, nil
}

// IsValid reports whether the role is a known role.
func IsValid(role Role) bool {
	_, ok := rolePermissions[role]
	return ok
}
