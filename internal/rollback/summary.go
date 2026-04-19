package rollback

import (
	"fmt"
	"io"
)

// WriteSummary writes a human-readable rollback summary to w.
func WriteSummary(w io.Writer, results []Result, dryRun bool) {
	mode := "applied"
	if dryRun {
		mode = "dry-run"
	}
	fmt.Fprintf(w, "Rollback summary [%s]\n", mode)
	fmt.Fprintf(w, "%-30s %-10s %s\n", "KEY", "STATUS", "DETAIL")

	for _, r := range results {
		status := "ok"
		detail := ""
		if r.Err != nil {
			status = "FAILED"
			detail = r.Err.Error()
		} else if dryRun {
			status = "would restore"
			detail = fmt.Sprintf("%q -> %q", r.OldVal, r.NewVal)
		}
		fmt.Fprintf(w, "%-30s %-10s %s\n", r.Key, status, detail)
	}

	failed := 0
	for _, r := range results {
		if r.Err != nil {
			failed++
		}
	}
	fmt.Fprintf(w, "\nTotal: %d  Failed: %d\n", len(results), failed)
}
