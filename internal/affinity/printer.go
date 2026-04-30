package affinity

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Print writes a formatted table of all affinity rules to w.
func Print(w io.Writer, s *Store) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tTARGET\tKIND")
	for _, r := range s.All() {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", r.Service, r.Target, r.Kind)
	}
	tw.Flush()
}
