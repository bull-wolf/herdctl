package priority

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Print writes a formatted priority table to w.
func Print(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no priority entries configured")
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Level != entries[j].Level {
			return entries[i].Level > entries[j].Level
		}
		return entries[i].Service < entries[j].Service
	})

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tPRIORITY")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%d\n", e.Service, e.Level)
	}
	_ = tw.Flush()
}
