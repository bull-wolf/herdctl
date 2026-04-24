package hooks

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Print writes a formatted table of registered hooks for all services to w.
func Print(r *Registry, services []string, w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tEVENT")

	sorted := make([]string, len(services))
	copy(sorted, services)
	sort.Strings(sorted)

	for _, svc := range sorted {
		events := r.List(svc)
		if len(events) == 0 {
			continue
		}
		evStrings := make([]string, len(events))
		for i, e := range events {
			evStrings[i] = string(e)
		}
		sort.Strings(evStrings)
		for _, ev := range evStrings {
			fmt.Fprintf(tw, "%s\t%s\n", svc, ev)
		}
	}
	tw.Flush()
}
