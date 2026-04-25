// Package importer provides functionality for importing secrets from
// external sources (JSON or env format) into a configured provider.
package importer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/vaultop/internal/provider"
)

// Format represents the input file format.
type Format string

const (
	FormatJSON Format = "json"
	FormatEnv  Format = "env"
)

// IsValid returns true if the format is supported.
func (f Format) IsValid() bool {
	return f == FormatJSON || f == FormatEnv
}

// Options configures an import operation.
type Options struct {
	Provider provider.Provider
	Format   Format
	DryRun   bool
}

// Import reads secrets from r and writes them to the provider.
// Returns a map of key->value that were imported.
func Import(r io.Reader, opts Options) (map[string]string, error) {
	if !opts.Format.IsValid() {
		return nil, fmt.Errorf("unsupported import format: %q", opts.Format)
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read input: %w", err)
	}

	var secrets map[string]string
	switch opts.Format {
	case FormatJSON:
		secrets, err = parseJSON(data)
	case FormatEnv:
		secrets, err = parseEnv(data)
	}
	if err != nil {
		return nil, err
	}

	if opts.DryRun {
		return secrets, nil
	}

	for k, v := range secrets {
		if err := opts.Provider.Set(k, v); err != nil {
			return nil, fmt.Errorf("set %q: %w", k, err)
		}
	}
	return secrets, nil
}

// ImportFile opens the file at path and calls Import.
// If opts.Format is empty, the format is inferred from the file extension
// (.json -> JSON, .env -> env).
func ImportFile(path string, opts Options) (map[string]string, error) {
	if opts.Format == "" {
		opts.Format = inferFormat(path)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	defer f.Close()
	return Import(f, opts)
}

// inferFormat guesses the Format from a file path's extension.
func inferFormat(path string) Format {
	switch {
	case strings.HasSuffix(path, ".json"):
		return FormatJSON
	case strings.HasSuffix(path, ".env"):
		return FormatEnv
	default:
		return ""
	}
}

func parseJSON(data []byte) (map[string]string, error) {
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse json: %w", err)
	}
	return m, nil
}

func parseEnv(data []byte) (map[string]string, error) {
	m := make(map[string]string)
	for lineNum, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid env line %d: %q", lineNum+1, line)
		}
		m[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return m, nil
}
