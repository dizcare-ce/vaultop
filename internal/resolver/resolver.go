// Package resolver provides secret key resolution with alias and fallback support.
package resolver

import (
	"errors"
	"fmt"
)

// ErrNotResolved is returned when a key cannot be resolved through any alias or fallback.
var ErrNotResolved = errors.New("resolver: key could not be resolved")

// Provider is the minimal interface required to fetch secrets.
type Provider interface {
	Get(key string) (string, error)
}

// Config holds resolver configuration.
type Config struct {
	// Aliases maps logical names to canonical provider keys.
	Aliases map[string]string
	// Fallbacks lists keys to try in order when the primary key is missing.
	Fallbacks map[string][]string
}

// Resolver resolves secret keys with alias expansion and ordered fallbacks.
type Resolver struct {
	provider  Provider
	aliases   map[string]string
	fallbacks map[string][]string
}

// New creates a Resolver backed by the given provider and config.
func New(p Provider, cfg Config) *Resolver {
	aliases := cfg.Aliases
	if aliases == nil {
		aliases = map[string]string{}
	}
	fallbacks := cfg.Fallbacks
	if fallbacks == nil {
		fallbacks = map[string][]string{}
	}
	return &Resolver{provider: p, aliases: aliases, fallbacks: fallbacks}
}

// Resolve returns the secret value for key, expanding aliases and trying
// fallbacks when the primary key is absent.
func (r *Resolver) Resolve(key string) (string, error) {
	candidates := r.candidates(key)
	for _, k := range candidates {
		val, err := r.provider.Get(k)
		if err == nil {
			return val, nil
		}
	}
	return "", fmt.Errorf("%w: %s", ErrNotResolved, key)
}

// candidates returns the ordered list of keys to try for the given logical key.
func (r *Resolver) candidates(key string) []string {
	primary := key
	if alias, ok := r.aliases[key]; ok {
		primary = alias
	}
	keys := []string{primary}
	if fb, ok := r.fallbacks[key]; ok {
		keys = append(keys, fb...)
	}
	return keys
}
