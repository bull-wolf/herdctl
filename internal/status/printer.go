package status

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// stateColor maps states to ANSI colour codes.
var stateColor = map[State]string{
	StateRunning:  "\033[32m", // green
	StateStarting: "\033[33m", // yellow
	StateStopped:  "\033[90m", // grey
	StateFailed:   "\033[31m", // red
}

const reset = "\033[0m"

// Print writes a formatted status table for all entries to w.
func Print(w io.Writer, entries []Entry) {
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Service < sorted[j].Service
	})

	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tSTATE\tPID\tSINCE")
	for _, e := range sorted {
		color := stateColor[e.State]
		pid := fmt.Sprintf("%d", e.PID)
		if e.PID == 0 {
			pid = "-"
		}
		fmt.Fprintf(tw, "%s\t%s%s%s\t%s\t%s\n",
			e.Service,
			color, e.State, reset,
			pid,
			e.UpdatedAt.Format("15:04:05"),
		)
	}
	tw.Flush()
}
