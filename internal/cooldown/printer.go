package cooldown

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Print writes a formatted cooldown status table to w.
func Print(w io.Writer, m *Manager) {
	entries := m.All()
	if len(entries) == 0 {
		fmt.Fprintln(w, "no cooldown entries")
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Service < entries[j].Service
	})

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tDURATION\tACTIVE\tREMAINING")
	for _, e := range entries {
		active := "no"
		if e.IsActive() {
			active = "yes"
		}
		remaining := "-"
		if e.IsActive() {
			remaining = e.RemainingTime().Round(1 * 1000000).String()
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", e.Service, e.Duration, active, remaining)
	}
	tw.Flush()
}
