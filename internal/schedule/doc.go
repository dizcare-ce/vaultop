// Package schedule implements rotation scheduling policies for vaultop secrets.
//
// A Policy defines how frequently a secret should be rotated via IntervalDays
// and tracks the LastRotated timestamp. CheckAll evaluates a slice of
// SecretSchedule entries and returns DueResult values indicating which secrets
// are overdue for rotation. FilterDue can then be used to extract only those
// secrets that require immediate action.
//
// Example usage:
//
//	policies := []schedule.SecretSchedule{
//		{Key: "db/password", Policy: schedule.Policy{IntervalDays: 7, LastRotated: lastTime}},
//	}
//	results, err := schedule.CheckAll(policies)
//	due := schedule.FilterDue(results)
package schedule
