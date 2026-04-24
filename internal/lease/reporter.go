package lease

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
	"time"
)

// ReportOptions controls how the lease report is rendered.
type ReportOptions struct {
	// Now is the reference time used to compute remaining TTL.
	// Defaults to time.Now() when zero.
	Now time.Time
	// ShowExpired includes expired leases in the output.
	ShowExpired bool
}

// Report writes a human-readable table of leases to w.
func Report(w io.Writer, leases []Lease, opts ReportOptions) error {
	now := opts.Now
	if now.IsZero() {
		now = time.Now()
	}

	if opts.ShowExpired == false {
		filtered := leases[:0]
		for _, l := range leases {
			if !l.IsExpired(now) {
				filtered = append(filtered, l)
			}
		}
		leases = filtered
	}

	sort.Slice(leases, func(i, j int) bool {
		return leases[i].Key < leases[j].Key
	})

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "KEY\tISSUED\tEXPIRES\tRENEWALS\tSTATUS")
	for _, l := range leases {
		status := "active"
		if l.IsExpired(now) {
			status = "expired"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t%s\n",
			l.Key,
			l.IssuedAt.Format(time.RFC3339),
			l.ExpiresAt().Format(time.RFC3339),
			l.Renewals,
			status,
		)
	}
	return tw.Flush()
}
