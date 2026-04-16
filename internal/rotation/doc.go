// Package rotation provides secret rotation logic for vaultop.
//
// A Rotator is constructed with a provider.Provider and Options, then
// Rotate is called with a list of secret IDs. For each ID a new value is
// produced by a ValueGenerator (defaulting to a 32-byte random base64
// string) and written back through the provider unless DryRun is enabled.
//
// Example:
//
//	p, _ := provider.New(provider.TypeStub, nil)
//	r := rotation.New(p, rotation.Options{})
//	results := r.Rotate(ctx, []string{"db/password", "api/key"})
//	for _, res := range results {
//		if res.Err != nil {
//			log.Printf("failed to rotate %s: %v", res.SecretID, res.Err)
//		}
//	}
package rotation
