package drain

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Print writes a human-readable table of active in-flight counts to w.
func Print(d *Drainer, services []string, w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tIN-FLIGHT")

	sorted := make([]string, len(services))
	copy(sorted, services)
	sort.Strings(sorted)

	for _, svc := range sorted {
		active := d.Active(svc)
		status := fmt.Sprintf("%d", active)
		if active == 0 {
			status = "idle"
		}
		fmt.Fprintf(tw, "%s\t%s\n", svc, status)
	}
	tw.Flush()
}
