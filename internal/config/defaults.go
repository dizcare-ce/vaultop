package config

import "time"

const (
	DefaultConfigFile        = "vaultop.yaml"
	DefaultRotationInterval  = 24 * time.Hour
	DefaultVersion           = "1"
)

// ApplyDefaults fills in zero-value fields with sensible defaults.
func (c *Config) ApplyDefaults() {
	if c.Version == "" {
		c.Version = DefaultVersion
	}
	if c.Defaults.RotationInterval == "" {
		c.Defaults.RotationInterval = DefaultRotationInterval.String()
	}
}

// RotationDuration parses the rotation interval string into a time.Duration.
// Falls back to DefaultRotationInterval on parse error.
func (c *Config) RotationDuration() time.Duration {
	d, err := time.ParseDuration(c.Defaults.RotationInterval)
	if err != nil {
		return DefaultRotationInterval
	}
	return d
}
