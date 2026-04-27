package backoff

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Print writes a formatted table of all backoff entries to w.
func Print(m *Manager, w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tSTRATEGY\tATTEMPTS\tLAST WAIT\tMAX WAIT")
	m.mu.Lock()
	entries := make([]*Entry, 0, len(m.entries))
	for _, e := range m.entries {
		entries = append(entries, e)
	}
	m.mu.Unlock()
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\t%s\n",
			e.Service,
			e.Strategy,
			e.Attempts,
			fmtDuration(e.LastWait),
			fmtDuration(e.MaxWait),
		)
	}
	tw.Flush()
}

func fmtDuration(d interface{ String() string }) string {
	if d == nil {
		return "-"
	}
	return d.String()
}
