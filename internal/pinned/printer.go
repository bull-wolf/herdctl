package pinned

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Print writes a formatted table of pinned service versions to w.
func Print(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no pinned services")
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Service < entries[j].Service
	})

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tVERSION\tPINNED BY\tPINNED AT\tREASON")
	for _, e := range entries {
		reason := e.Reason
		if reason == "" {
			reason = "-"
		}
		pinnedBy := e.PinnedBy
		if pinnedBy == "" {
			pinnedBy = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n",
			e.Service,
			e.Version,
			pinnedBy,
			e.PinnedAt.Format("2006-01-02 15:04:05"),
			reason,
		)
	}
	_ = tw.Flush()
}
