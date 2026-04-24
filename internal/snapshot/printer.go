package snapshot

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Print writes a formatted table of all snapshots to w.
func Print(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no snapshots captured")
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Service < entries[j].Service
	})

	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tSTATUS\tHEALTH\tCAPTURED AT")
	for _, e := range entries {
		health := e.Health
		if health == "" {
			health = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			e.Service,
			e.Status,
			health,
			e.CapturedAt.Format("2006-01-02 15:04:05"),
		)
	}
	tw.Flush()
}
