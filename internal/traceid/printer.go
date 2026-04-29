package traceid

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Print writes a formatted table of trace IDs to w.
func Print(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no trace IDs recorded")
		return
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Service < entries[j].Service
	})
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tTRACE ID\tCREATED AT")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%s\n",
			e.Service,
			e.TraceID,
			e.CreatedAt.Format("15:04:05.000"),
		)
	}
	tw.Flush()
}
