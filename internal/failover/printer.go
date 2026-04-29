package failover

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Print writes a formatted table of failover entries to w.
func Print(w io.Writer, entries []Entry) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tSTRATEGY\tACTIVE\tCURRENT\tTARGETS")
	for _, e := range entries {
		activeStr := "no"
		if e.Active {
			activeStr = "yes"
		}
		current := e.Current
		if current == "" {
			current = "-"
		}
		targets := formatTargets(e.Targets)
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n",
			e.Service, e.Strategy, activeStr, current, targets)
	}
	tw.Flush()
}

func formatTargets(targets []string) string {
	if len(targets) == 0 {
		return "-"
	}
	out := ""
	for i, t := range targets {
		if i > 0 {
			out += ","
		}
		out += t
	}
	return out
}
