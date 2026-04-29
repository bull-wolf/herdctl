package eventbus

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Print writes a formatted table of active subscriptions to w.
func Print(w io.Writer, b *Bus) {
	entries := b.List()
	if len(entries) == 0 {
		fmt.Fprintln(w, "no active eventbus subscriptions")
		return
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Service != entries[j].Service {
			return entries[i].Service < entries[j].Service
		}
		return entries[i].Event < entries[j].Event
	})
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tEVENT")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\n", e.Service, e.Event)
	}
	_ = tw.Flush()
}
