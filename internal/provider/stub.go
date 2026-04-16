package provider

import "fmt"

// stubProvider is a minimal in-memory provider used as a placeholder
// until real provider implementations are added.
type stubProvider struct {
	kind Type
	opts map[string]string
	store map[string]string
}

func (s *stubProvider) GetSecret(name string) (string, error) {
	v, ok := s.store[name]
	if !ok {
		return "", fmt.Errorf("%s: secret %q not found", s.kind, name)
	}
	return v, nil
}

func (s *stubProvider) SetSecret(name, value string) error {
	s.store[name] = value
	return nil
}

func (s *stubProvider) DeleteSecret(name string) error {
	if _, ok := s.store[name]; !ok {
		return fmt.Errorf("%s: secret %q not found", s.kind, name)
	}
	delete(s.store, name)
	return nil
}

func (s *stubProvider) ListSecrets(prefix string) ([]string, error) {
	var keys []string
	for k := range s.store {
		if prefix == "" || len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

func newStub(kind Type, opts map[string]string) *stubProvider {
	return &stubProvider{kind: kind, opts: opts, store: make(map[string]string)}
}

func newAWSProvider(opts map[string]string) (Provider, error)   { return newStub(AWS, opts), nil }
func newGCPProvider(opts map[string]string) (Provider, error)   { return newStub(GCP, opts), nil }
func newAzureProvider(opts map[string]string) (Provider, error) { return newStub(Azure, opts), nil }
func newVaultProvider(opts map[string]string) (Provider, error) { return newStub(Vault, opts), nil }
