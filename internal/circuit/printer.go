package circuit

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Print writes a formatted table of all circuit breaker entries to w.
func Print(w io.Writer, entries []Entry) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tSTATE\tFAILURES\tOPENED AT\tLAST ERROR")
	for _, e := range entries {
		opened := "-"
		if e.State == StateOpen || e.State == StateHalfOpen {
			opened = e.OpenedAt.Format(time.RFC3339)
		}
		lastErr := "-"
		if e.LastError != nil {
			lastErr = e.LastError.Error()
		}
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\t%s\n",
			e.Service, e.State, e.Failures, opened, lastErr)
	}
	tw.Flush()
}
