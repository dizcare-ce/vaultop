// Package backup provides functionality for creating and restoring
// versioned backups of secrets before rotation or bulk changes.
package backup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single backup file record.
type Entry struct {
	ID        string            `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	Secrets   map[string]string `json:"secrets"`
}

// Manager handles backup storage in a directory.
type Manager struct {
	dir string
}

// NewManager returns a Manager that stores backups under dir.
func NewManager(dir string) (*Manager, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("backup: create dir: %w", err)
	}
	return &Manager{dir: dir}, nil
}

// Save writes secrets to a timestamped backup file and returns the entry ID.
func (m *Manager) Save(secrets map[string]string) (string, error) {
	id := time.Now().UTC().Format("20060102T150405Z")
	entry := Entry{
		ID:        id,
		CreatedAt: time.Now().UTC(),
		Secrets:   secrets,
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return "", fmt.Errorf("backup: marshal: %w", err)
	}
	path := filepath.Join(m.dir, id+".json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return "", fmt.Errorf("backup: write: %w", err)
	}
	return id, nil
}

// Load reads a backup entry by ID.
func (m *Manager) Load(id string) (*Entry, error) {
	path := filepath.Join(m.dir, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("backup: read %s: %w", id, err)
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("backup: unmarshal: %w", err)
	}
	return &entry, nil
}

// List returns all backup IDs sorted by filename (chronological).
func (m *Manager) List() ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(m.dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("backup: list: %w", err)
	}
	ids := make([]string, 0, len(matches))
	for _, p := range matches {
		base := filepath.Base(p)
		ids = append(ids, base[:len(base)-5])
	}
	return ids, nil
}
