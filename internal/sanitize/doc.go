// Package sanitize provides helpers for cleaning and normalising secret keys
// and values before they are stored in or retrieved from a provider.
//
// # Value sanitisation
//
// Apply and ApplyMap clean individual secret values according to an Options
// struct. The default options trim surrounding whitespace, strip non-printable
// control characters, and reject empty results.
//
// # Key normalisation
//
// NormaliseKey converts an arbitrary key string to a canonical lowercase,
// underscore-separated form so that keys coming from different sources (e.g.
// environment files vs. JSON exports) can be compared reliably.
package sanitize
