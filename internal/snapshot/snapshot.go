// Package snapshot provides functionality to capture and restore
// the current state of secrets from a provider.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/vaultop/internal/provider"
)

// Snapshot holds a point-in-time capture of secrets.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Provider  string            `json:"provider"`
	Secrets   map[string]string `json:"secrets"`
}

// Take captures all secrets from the given keys using the provider.
func Take(p provider.Provider, keys []string) (*Snapshot, error) {
	secrets := make(map[string]string, len(keys))
	for _, k := range keys {
		val, err := p.Get(k)
		if err != nil {
			return nil, fmt.Errorf("snapshot: get %q: %w", k, err)
		}
		secrets[k] = val
	}
	return &Snapshot{
		CreatedAt: time.Now().UTC(),
		Provider:  string(p.Type()),
		Secrets:   secrets,
	}, nil
}

// Save writes the snapshot as JSON to the given file path.
func Save(s *Snapshot, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// Load reads a snapshot from a JSON file.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &s, nil
}

// Restore writes all secrets from the snapshot back into the provider.
func Restore(p provider.Provider, s *Snapshot) error {
	for k, v := range s.Secrets {
		if err := p.Set(k, v); err != nil {
			return fmt.Errorf("snapshot: restore %q: %w", k, err)
		}
	}
	return nil
}
