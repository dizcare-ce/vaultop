// Package redact provides utilities for masking secret values in logs and output.
package redact

import "strings"

// Mode controls how a value is redacted.
type Mode string

const (
	ModeFull    Mode = "full"    // replace entire value with mask
	ModePartial Mode = "partial" // show last N chars
	ModeHash    Mode = "hash"    // show length only
)

// DefaultMask is the string used to replace secret values.
const DefaultMask = "***"

// Options configures redaction behaviour.
type Options struct {
	Mode        Mode
	ShowSuffix  int    // chars to reveal at end (ModePartial only)
	Mask        string // defaults to DefaultMask
}

func (o Options) mask() string {
	if o.Mask == "" {
		return DefaultMask
	}
	return o.Mask
}

// Redact returns a redacted representation of value.
func Redact(value string, opts Options) string {
	switch opts.Mode {
	case ModePartial:
		n := opts.ShowSuffix
		if n <= 0 {
			n = 4
		}
		if len(value) <= n {
			return opts.mask()
		}
		return opts.mask() + value[len(value)-n:]
	case ModeHash:
		return strings.Repeat("*", len(value))
	default: // ModeFull
		return opts.mask()
	}
}

// RedactMap returns a copy of m with all values redacted.
func RedactMap(m map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = Redact(v, opts)
	}
	return out
}
