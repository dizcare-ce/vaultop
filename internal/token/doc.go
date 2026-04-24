// Package token provides short-lived, HMAC-signed tokens for
// service-to-service authentication within vaultop.
//
// # Overview
//
// A Manager issues tokens for a named subject and validates them on
// subsequent requests. Each token embeds an expiry timestamp and a
// random nonce so replayed or tampered values are rejected.
//
// A Store complements the Manager by tracking issued tokens and
// supporting explicit revocation before natural expiry. Calling
// Purge periodically reclaims memory for tokens that have expired.
//
// # Usage
//
//	mgr, _ := token.New(secret, 15*time.Minute)
//	tok, _ := mgr.Issue("rotation-worker")
//
//	subject, err := mgr.Validate(tok.Value)
//
// Tokens should be transmitted over TLS only; the HMAC prevents
// forgery but does not encrypt the payload.
package token
