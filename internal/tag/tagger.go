// Package tag provides utilities for reading and writing metadata tags
// on secrets stored in a provider.
package tag

import (
	"context"
	"fmt"
	"sort"
)

// Provider is the subset of the secret provider interface needed for tagging.
type Provider interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	List(ctx context.Context) ([]string, error)
}

// tagKey returns the internal key used to store tags for a secret.
func tagKey(secret string) string {
	return fmt.Sprintf("__tags__/%s", secret)
}

// Set writes a single tag (name=value) for the given secret key.
func Set(ctx context.Context, p Provider, secret, name, value string) error {
	tags, err := GetAll(ctx, p, secret)
	if err != nil {
		return err
	}
	tags[name] = value
	return persist(ctx, p, secret, tags)
}

// Get returns the value of a single tag for a secret.
func Get(ctx context.Context, p Provider, secret, name string) (string, bool, error) {
	tags, err := GetAll(ctx, p, secret)
	if err != nil {
		return "", false, err
	}
	v, ok := tags[name]
	return v, ok, nil
}

// Delete removes a tag from a secret.
func Delete(ctx context.Context, p Provider, secret, name string) error {
	tags, err := GetAll(ctx, p, secret)
	if err != nil {
		return err
	}
	delete(tags, name)
	return persist(ctx, p, secret, tags)
}

// GetAll returns all tags for a secret as a map.
func GetAll(ctx context.Context, p Provider, secret string) (map[string]string, error) {
	raw, err := p.Get(ctx, tagKey(secret))
	if err != nil {
		// No tags stored yet — return empty map.
		return map[string]string{}, nil
	}
	return parse(raw)
}

// ListTagged returns all secret keys that carry at least one tag.
func ListTagged(ctx context.Context, p Provider) ([]string, error) {
	keys, err := p.List(ctx)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, k := range keys {
		if len(k) > 9 && k[:9] == "__tags__/" {
			out = append(out, k[9:])
		}
	}
	sort.Strings(out)
	return out, nil
}

func persist(ctx context.Context, p Provider, secret string, tags map[string]string) error {
	return p.Set(ctx, tagKey(secret), encode(tags))
}
