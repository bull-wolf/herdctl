package pipeline

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"
)

var statusIcon = map[string]string{
	"pending": "○",
	"running": "►",
	"done":    "✓",
	"failed":  "✗",
}

// Print writes a formatted pipeline table for all stages of a service.
func Print(w io.Writer, service string, stages []*Stage) {
	fmt.Fprintf(w, "Pipeline: %s\n", service)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "  STAGE\tSTATUS\tDURATION")
	fmt.Fprintln(tw, "  "+strings.Repeat("-", 36))
	for _, s := range stages {
		icon := statusIcon[s.Status]
		if icon == "" {
			icon = "?"
		}
		dur := "-"
		if !s.Started.IsZero() {
			end := s.Finished
			if end.IsZero() {
				end = time.Now()
			}
			dur = end.Sub(s.Started).Round(time.Millisecond).String()
		}
		fmt.Fprintf(tw, "  %s %s\t%s\t%s\n", icon, s.Name, s.Status, dur)
	}
	tw.Flush()
}
