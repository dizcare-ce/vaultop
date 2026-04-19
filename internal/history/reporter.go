package history

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Report writes a human-readable summary of rotation history to w.
// Only entries for the provided keys are included; pass nil to include all.
// Only entries at or after since are included; pass zero time to include all.
func Report(w io.Writer, s *Store, keys []string, since time.Time) error {
	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[k] = true
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SECRET KEY\tPROVIDER\tROTATED AT\tSTATUS")
	fmt.Fprintln(tw, "----------\t--------\t----------\t------")

	var count int
	for _, e := range s.All() {
		if len(keySet) > 0 && !keySet[e.SecretKey] {
			continue
		}
		if !since.IsZero() && e.RotatedAt.Before(since) {
			continue
		}
		status := "ok"
		if !e.Success {
			status = "FAILED"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			e.SecretKey,
			e.Provider,
			e.RotatedAt.Format(time.RFC3339),
			status,
		)
		count++
	}
	if count == 0 {
		fmt.Fprintln(tw, "(no entries)")
	}
	if err := tw.Flush(); err != nil {
		return fmt.Errorf("flushing report output: %w", err)
	}
	return nil
}
