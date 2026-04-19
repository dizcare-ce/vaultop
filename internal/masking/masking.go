// Package masking provides configurable secret masking for log output and display.
package masking

import "strings"

// Mode controls how a secret value is masked.
type Mode string

const (
	ModeFull    Mode = "full"    // replace entire value with mask
	ModePartial Mode = "partial" // show last N chars
	ModePrefix  Mode = "prefix"  // show first N chars
)

// Config holds masking options.
type Config struct {
	Mode       Mode
	Mask       string
	VisibleLen int // used by partial/prefix modes
}

// DefaultConfig returns sensible masking defaults.
func DefaultConfig() Config {
	return Config{
		Mode:       ModeFull,
		Mask:       "***",
		VisibleLen: 4,
	}
}

// Apply masks value according to cfg.
func Apply(value string, cfg Config) string {
	if value == "" {
		return cfg.Mask
	}
	switch cfg.Mode {
	case ModePartial:
		if len(value) <= cfg.VisibleLen {
			return cfg.Mask
		}
		return cfg.Mask + value[len(value)-cfg.VisibleLen:]
	case ModePrefix:
		if len(value) <= cfg.VisibleLen {
			return cfg.Mask
		}
		return value[:cfg.VisibleLen] + cfg.Mask
	default:
		return cfg.Mask
	}
}

// ApplyMap masks all values in a map, returning a new map.
func ApplyMap(secrets map[string]string, cfg Config) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = Apply(v, cfg)
	}
	return out
}

// ContainsSensitive returns true if s contains any of the provided secret values.
func ContainsSensitive(s string, secrets []string) bool {
	for _, sec := range secrets {
		if sec != "" && strings.Contains(s, sec) {
			return true
		}
	}
	return false
}
