package labels

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Print writes a formatted table of labels for a service to w.
// If the service has no labels an informational message is printed.
func Print(w io.Writer, service string, s *Store) {
	labels := s.All(service)
	if len(labels) == 0 {
		fmt.Fprintf(w, "no labels for service %q\n", service)
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "KEY\tVALUE")

	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(tw, "%s\t%s\n", k, labels[k])
	}
	tw.Flush()
}
