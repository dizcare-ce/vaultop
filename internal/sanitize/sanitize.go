// Package sanitize provides utilities for cleaning and normalising secret
// values before they are written to a provider or returned to callers.
package sanitize

import (
	"errors"
	"strings"
	"unicode"
)

// ErrEmptyValue is returned when a value is empty after sanitisation.
var ErrEmptyValue = errors.New("sanitize: value is empty after sanitisation")

// Options controls which sanitisation steps are applied.
type Options struct {
	// TrimSpace removes leading and trailing whitespace.
	TrimSpace bool
	// StripControl removes non-printable control characters.
	StripControl bool
	// MaxLength truncates the value to at most this many runes. 0 means no limit.
	MaxLength int
	// RejectEmpty causes Apply to return ErrEmptyValue when the result is empty.
	RejectEmpty bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		TrimSpace:    true,
		StripControl: true,
		MaxLength:    0,
		RejectEmpty:  true,
	}
}

// Apply sanitises a single secret value according to opts.
func Apply(value string, opts Options) (string, error) {
	if opts.TrimSpace {
		value = strings.TrimSpace(value)
	}
	if opts.StripControl {
		value = stripControl(value)
	}
	if opts.MaxLength > 0 {
		runes := []rune(value)
		if len(runes) > opts.MaxLength {
			value = string(runes[:opts.MaxLength])
		}
	}
	if opts.RejectEmpty && value == "" {
		return "", ErrEmptyValue
	}
	return value, nil
}

// ApplyMap sanitises every value in m, returning a new map.
// Keys whose values fail sanitisation are omitted and their errors collected.
func ApplyMap(m map[string]string, opts Options) (map[string]string, map[string]error) {
	out := make(map[string]string, len(m))
	errs := make(map[string]error)
	for k, v := range m {
		clean, err := Apply(v, opts)
		if err != nil {
			errs[k] = err
			continue
		}
		out[k] = clean
	}
	return out, errs
}

func stripControl(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\t' && r != '\n' {
			return -1
		}
		return r
	}, s)
}
