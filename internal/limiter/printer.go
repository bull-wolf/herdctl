package limiter

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Print writes a formatted table of limiter entries to w.
func Print(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no concurrency limits configured")
		return
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Service < entries[j].Service
	})
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tACTIVE\tMAX\tUTILIZATION")
	for _, e := range entries {
		utilPct := 0
		if e.Max > 0 {
			utilPct = (e.Active * 100) / e.Max
		}
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d%%\n", e.Service, e.Active, e.Max, utilPct)
	}
	_ = tw.Flush()
}
