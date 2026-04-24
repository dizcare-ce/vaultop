package sanitize

import (
	"strings"
	"unicode"
)

// NormaliseKey returns a canonical form of a secret key: lowercase, with
// runs of non-alphanumeric characters replaced by a single underscore, and
// leading/trailing underscores removed.
func NormaliseKey(key string) string {
	key = strings.ToLower(strings.TrimSpace(key))
	var b strings.Builder
	prevUnderscore := false
	for _, r := range key {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			prevUnderscore = false
		} else {
			if !prevUnderscore {
				b.WriteRune('_')
				prevUnderscore = true
			}
		}
	}
	return strings.Trim(b.String(), "_")
}

// NormaliseMap returns a new map with all keys normalised via NormaliseKey.
// If two keys normalise to the same value, the last one wins (iteration order
// is undefined, so callers should avoid duplicate normalised keys).
func NormaliseMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[NormaliseKey(k)] = v
	}
	return out
}
