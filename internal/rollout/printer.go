package rollout

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Print writes a formatted rollout status table to w.
func Print(w io.Writer, states []*State) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tSTRATEGY\tSTATUS\tAGE")
	for _, s := range states {
		status := statusLabel(s)
		age := time.Since(s.StartedAt).Round(time.Second)
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", s.Service, s.Strategy, status, age)
	}
	tw.Flush()
}

func statusLabel(s *State) string {
	if !s.Done {
		return "in-progress"
	}
	if s.Err != nil {
		return fmt.Sprintf("failed: %s", s.Err.Error())
	}
	return "done"
}
