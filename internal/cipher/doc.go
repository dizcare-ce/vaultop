// Package cipher provides AES-256-GCM encryption and decryption for secret
// values stored by vaultop.
//
// # Usage
//
//	c, err := cipher.New(key, cipher.AES256GCM)
//	ct, err := c.Encrypt([]byte("secret"))
//	pt, err := c.Decrypt(ct)
//
// EncryptedStore wraps any SecretStore to transparently encrypt values at rest:
//
//	store := cipher.NewEncryptedStore(provider, c)
//	_ = store.Set(ctx, "db/password", "s3cr3t")
//	v, _ := store.Get(ctx, "db/password") // returns "s3cr3t"
package cipher
