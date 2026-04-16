# vaultop

> CLI tool for managing and rotating secrets across multiple cloud providers

## Installation

```bash
go install github.com/yourusername/vaultop@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultop/releases).

## Usage

```bash
# List all secrets in a provider
vaultop list --provider aws

# Rotate a secret
vaultop rotate --provider gcp --secret my-api-key

# Sync secrets across providers
vaultop sync --from aws --to azure --secret db-password
```

### Supported Providers

- AWS Secrets Manager
- Google Cloud Secret Manager
- Azure Key Vault

### Configuration

```bash
# Initialize vaultop with your provider credentials
vaultop init

# Configuration is stored in ~/.vaultop/config.yaml
```

## Example

```bash
$ vaultop rotate --provider aws --secret prod/db-password
✔ Fetching current secret...
✔ Generating new secret value...
✔ Updating secret in AWS Secrets Manager...
✔ Secret rotated successfully.
```

## License

[MIT](LICENSE)