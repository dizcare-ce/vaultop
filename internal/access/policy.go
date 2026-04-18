package access

import "fmt"

// Entry binds an identity (user or service account) to a Role.
type Entry struct {
	Identity string `yaml:"identity" json:"identity"`
	Role     Role   `yaml:"role"     json:"role"`
}

// Policy holds a collection of access entries.
type Policy struct {
	Entries []Entry `yaml:"entries" json:"entries"`
}

// Lookup returns the Role for the given identity, or an error if not found.
func (p *Policy) Lookup(identity string) (Role, error) {
	for _, e := range p.Entries {
		if e.Identity == identity {
			return e.Role, nil
		}
	}
	return "", fmt.Errorf("access: identity %q not found in policy", identity)
}

// Allow returns nil when identity is permitted to perform op.
func (p *Policy) Allow(identity string, op Op) error {
	role, err := p.Lookup(identity)
	if err != nil {
		return err
	}
	return Check(role, op)
}

// Validate checks that all entries reference known roles.
func (p *Policy) Validate() error {
	for _, e := range p.Entries {
		if e.Identity == "" {
			return fmt.Errorf("access: entry has empty identity")
		}
		if !IsValid(e.Role) {
			return fmt.Errorf("access: unknown role %q for identity %q", e.Role, e.Identity)
		}
	}
	return nil
}
