package fence

import (
	"context"
	"fmt"
)

// WriteFunc is the signature of a guarded write operation.
type WriteFunc func(ctx context.Context, key, value string) error

// GuardedWriter wraps a WriteFunc and enforces fencing token validation
// before each write is allowed through.
type GuardedWriter struct {
	manager  *Manager
	resource string
	inner    WriteFunc
}

// NewGuardedWriter returns a GuardedWriter for the given resource.
func NewGuardedWriter(m *Manager, resource string, inner WriteFunc) *GuardedWriter {
	return &GuardedWriter{manager: m, resource: resource, inner: inner}
}

// Write performs the underlying write only if tok is still the current
// fencing token for the resource. Returns ErrStaleFence otherwise.
func (g *GuardedWriter) Write(ctx context.Context, tok Token, key, value string) error {
	if err := g.manager.Check(g.resource, tok); err != nil {
		return fmt.Errorf("guarded write rejected for key %q: %w", key, err)
	}
	return g.inner(ctx, key, value)
}
