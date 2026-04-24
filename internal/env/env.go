// Package env provides utilities for expanding secret references embedded
// in environment-variable style strings such as "${secret:db/password}".
// It resolves each placeholder through the configured provider and returns
// the substituted result.
package env

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// Provider is the minimal interface required to look up a secret value.
type Provider interface {
	Get(ctx context.Context, key string) (string, error)
}

// placeholder matches "${secret:some/key}" tokens inside a string.
var placeholder = regexp.MustCompile(`\$\{secret:([^}]+)\}`)

// Expand replaces every "${secret:<key>}" token in s with the value
// returned by the provider. The first resolution error aborts expansion
// and is returned to the caller.
func Expand(ctx context.Context, p Provider, s string) (string, error) {
	var expandErr error

	result := placeholder.ReplaceAllStringFunc(s, func(match string) string {
		if expandErr != nil {
			return match
		}
		sub := placeholder.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		key := strings.TrimSpace(sub[1])
		val, err := p.Get(ctx, key)
		if err != nil {
			expandErr = fmt.Errorf("env: resolve %q: %w", key, err)
			return match
		}
		return val
	})

	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}

// ExpandMap calls Expand on every value in m and returns a new map with
// substituted values. The original map is not modified.
func ExpandMap(ctx context.Context, p Provider, m map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(m))
	for k, v := range m {
		expanded, err := Expand(ctx, p, v)
		if err != nil {
			return nil, fmt.Errorf("env: key %q: %w", k, err)
		}
		out[k] = expanded
	}
	return out, nil
}
