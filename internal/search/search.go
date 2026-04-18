// Package search provides filtering and lookup utilities for secrets
// stored across a provider.
package search

import (
	"context"
	"strings"
)

// Provider is the subset of provider.Provider needed for searching.
type Provider interface {
	List(ctx context.Context) ([]string, error)
	Get(ctx context.Context, key string) (string, error)
}

// Result holds a single matched secret.
type Result struct {
	Key   string
	Value string
}

// Options controls search behaviour.
type Options struct {
	// Prefix filters keys that start with this string.
	Prefix string
	// Contains filters keys that contain this substring.
	Contains string
	// ValueContains filters results whose value contains this substring.
	ValueContains string
}

// Find lists all secrets from the provider and returns those matching opts.
func Find(ctx context.Context, p Provider, opts Options) ([]Result, error) {
	keys, err := p.List(ctx)
	if err != nil {
		return nil, err
	}

	var results []Result
	for _, key := range keys {
		if opts.Prefix != "" && !strings.HasPrefix(key, opts.Prefix) {
			continue
		}
		if opts.Contains != "" && !strings.Contains(key, opts.Contains) {
			continue
		}

		if opts.ValueContains == "" {
			results = append(results, Result{Key: key})
			continue
		}

		val, err := p.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		if strings.Contains(val, opts.ValueContains) {
			results = append(results, Result{Key: key, Value: val})
		}
	}
	return results, nil
}
