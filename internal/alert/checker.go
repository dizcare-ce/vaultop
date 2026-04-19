package alert

import "time"

// ExpirySource returns the expiry time for a key, and whether it exists.
type ExpirySource interface {
	Expiry(key string) (time.Time, bool)
}

// CheckAll evaluates all rules against the expiry source at the given time
// and dispatches any triggered alerts to the notifier.
func CheckAll(rules []Rule, src ExpirySource, n Notifier, now time.Time) []error {
	var errs []error
	for _, rule := range rules {
		expiry, ok := src.Expiry(rule.Key)
		if !ok {
			continue
		}
		if a, fired := CheckExpiry(rule, expiry, now); fired {
			if err := n.Notify(a); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

// FilterFired returns only the alerts that would fire given the rules and source.
func FilterFired(rules []Rule, src ExpirySource, now time.Time) []Alert {
	var out []Alert
	for _, rule := range rules {
		expiry, ok := src.Expiry(rule.Key)
		if !ok {
			continue
		}
		if a, fired := CheckExpiry(rule, expiry, now); fired {
			out = append(out, a)
		}
	}
	return out
}
