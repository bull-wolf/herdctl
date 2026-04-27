package quota

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Print writes a formatted quota table to w.
func Print(w io.Writer, entries []Entry) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tKIND\tUSED\tLIMIT\tSTATUS")
	for _, e := range entries {
		status := "ok"
		if e.Exceeded() {
			status = "EXCEEDED"
		}
		fmt.Fprintf(tw, "%s\t%s\t%.2f\t%.2f\t%s\n",
			e.Service, e.Kind, e.Used, e.Limit, status)
	}
	tw.Flush()
}
