package tag

import (
	"fmt"
	"sort"
	"strings"
)

// encode serialises a tag map as "key=value" lines joined by newlines.
// Keys and values must not contain '=' or newlines.
func encode(tags map[string]string) string {
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	lines := make([]string, 0, len(keys))
	for _, k := range keys {
		lines = append(lines, fmt.Sprintf("%s=%s", k, tags[k]))
	}
	return strings.Join(lines, "\n")
}

// parse deserialises the output of encode back into a map.
func parse(raw string) (map[string]string, error) {
	out := map[string]string{}
	if raw == "" {
		return out, nil
	}
	for _, line := range strings.Split(raw, "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("tag: malformed entry %q", line)
		}
		out[parts[0]] = parts[1]
	}
	return out, nil
}
