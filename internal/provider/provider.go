package provider

import "fmt"

// Provider defines the interface for secret management across cloud providers.
type Provider interface {
	GetSecret(name string) (string, error)
	SetSecret(name, value string) error
	DeleteSecret(name string) error
	ListSecrets(prefix string) ([]string, error)
}

// Type represents a supported cloud provider.
type Type string

const (
	AWS   Type = "aws"
	GCP   Type = "gcp"
	Azure Type = "azure"
	Vault Type = "vault"
)

// SupportedProviders lists all valid provider types.
var SupportedProviders = []Type{AWS, GCP, Azure, Vault}

// IsValid returns true if the provider type is supported.
func (t Type) IsValid() bool {
	for _, p := range SupportedProviders {
		if t == p {
			return true
		}
	}
	return false
}

// New returns a Provider implementation for the given type and options.
func New(t Type, opts map[string]string) (Provider, error) {
	switch t {
	case AWS:
		return newAWSProvider(opts)
	case GCP:
		return newGCPProvider(opts)
	case Azure:
		return newAzureProvider(opts)
	case Vault:
		return newVaultProvider(opts)
	default:
		return nil, fmt.Errorf("unsupported provider: %q", t)
	}
}
