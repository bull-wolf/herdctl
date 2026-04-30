package retention

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Print writes a formatted table of retention policies to w.
func Print(w io.Writer, policies []Policy) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tMAX AGE\tMAX ITEMS")
	for _, p := range policies {
		age := formatDuration(p.MaxAge)
		items := fmt.Sprintf("%d", p.MaxItems)
		if p.MaxItems == 0 {
			items = "unlimited"
		}
		if p.MaxAge == 0 {
			age = "unlimited"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\n", p.Service, age, items)
	}
	tw.Flush()
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "unlimited"
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	return fmt.Sprintf("%dh", int(d.Hours()))
}
