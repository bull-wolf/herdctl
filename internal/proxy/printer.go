package proxy

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Print writes a formatted table of all proxy entries to w.
func Print(p *Proxy, w io.Writer) {
	entries := p.All()
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Service < entries[j].Service
	})

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tPORT\tUPSTREAM\tSTATUS")
	for _, e := range entries {
		status := "enabled"
		if !e.Enabled {
			status = "disabled"
		}
		fmt.Fprintf(tw, "%s\t%d\t%s\t%s\n", e.Service, e.Port, e.Upstream, status)
	}
	tw.Flush()
}
