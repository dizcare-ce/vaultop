// Package export provides functionality for exporting secrets to various formats.
package export

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents an export output format.
type Format string

const (
	FormatJSON Format = "json"
	FormatEnv  Format = "env"
	FormatDotEnv Format = "dotenv"
)

// IsValid reports whether f is a recognised export format.
func (f Format) IsValid() bool {
	switch f {
	case FormatJSON, FormatEnv, FormatDotEnv:
		return true
	}
	return false
}

// Write serialises secrets into the requested format and writes to w.
func Write(w io.Writer, secrets map[string]string, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, secrets)
	case FormatEnv, FormatDotEnv:
		return writeEnv(w, secrets)
	default:
		return fmt.Errorf("export: unsupported format %q", format)
	}
}

func writeJSON(w io.Writer, secrets map[string]string) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(secrets)
}

func writeEnv(w io.Writer, secrets map[string]string) error {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(secrets[k])
		sb.WriteByte('\n')
	}
	_, err := fmt.Fprint(w, sb.String())
	return err
}
