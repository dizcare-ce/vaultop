package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Provider represents a supported cloud secret provider.
type Provider string

const (
	ProviderAWS   Provider = "aws"
	ProviderGCP   Provider = "gcp"
	ProviderAzure Provider = "azure"
	ProviderVault Provider = "vault"
)

// ProviderConfig holds provider-specific settings.
type ProviderConfig struct {
	Type    Provider          `yaml:"type"`
	Region  string            `yaml:"region,omitempty"`
	Project string            `yaml:"project,omitempty"`
	Options map[string]string `yaml:"options,omitempty"`
}

// Config is the top-level vaultop configuration.
type Config struct {
	Version   string                    `yaml:"version"`
	Providers map[string]ProviderConfig  `yaml:"providers"`
	Defaults  DefaultsConfig            `yaml:"defaults"`
}

// DefaultsConfig holds global default settings.
type DefaultsConfig struct {
	RotationInterval string `yaml:"rotation_interval"`
	DryRun           bool   `yaml:"dry_run"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// Validate checks that the configuration is semantically valid.
func (c *Config) Validate() error {
	if c.Version == "" {
		return fmt.Errorf("version is required")
	}
	if len(c.Providers) == 0 {
		return fmt.Errorf("at least one provider must be defined")
	}
	for name, p := range c.Providers {
		if name == "" {
			return fmt.Errorf("provider name must not be empty")
		}
		switch p.Type {
		case ProviderAWS, ProviderGCP, ProviderAzure, ProviderVault:
		default:
			return fmt.Errorf("provider %q has unknown type %q", name, p.Type)
		}
	}
	return nil
}
