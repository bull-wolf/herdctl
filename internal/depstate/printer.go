package depstate

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Print writes a human-readable table of dependency states to w.
func Print(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no dependency state entries")
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Service < entries[j].Service
	})

	fmt.Fprintf(w, "%-20s %-10s %s\n", "SERVICE", "STATE", "BLOCKING")
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 55))
	for _, e := range entries {
		blocking := "-"
		if len(e.Blocking) > 0 {
			blocking = strings.Join(e.Blocking, ", ")
		}
		fmt.Fprintf(w, "%-20s %-10s %s\n", e.Service, e.State, blocking)
	}
}
