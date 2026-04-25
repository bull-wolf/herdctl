package metrics

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Print writes a formatted metrics summary table to w.
// It shows the latest recorded sample for every known service.
func Print(w io.Writer, c *Collector) {
	svcs := c.Services()
	if len(svcs) == 0 {
		fmt.Fprintln(w, "no metrics recorded")
		return
	}
	sort.Strings(svcs)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tCPU%\tMEMORY (MB)\tUPTIME (s)\tSAMPLED AT")
	fmt.Fprintln(tw, "-------\t----\t-----------\t----------\t----------")

	for _, svc := range svcs {
		e, ok := c.Latest(svc)
		if !ok {
			continue
		}
		fmt.Fprintf(tw, "%s\t%.1f\t%.1f\t%d\t%s\n",
			e.Service,
			e.CPU,
			e.MemoryMB,
			e.UptimeSec,
			e.Timestamp.Format("15:04:05"),
		)
	}
	tw.Flush()
}

// PrintService writes a formatted metrics summary table to w for a single
// named service. If no metrics have been recorded for the service, a
// descriptive message is written instead.
func PrintService(w io.Writer, c *Collector, service string) {
	e, ok := c.Latest(service)
	if !ok {
		fmt.Fprintf(w, "no metrics recorded for service %q\n", service)
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tCPU%\tMEMORY (MB)\tUPTIME (s)\tSAMPLED AT")
	fmt.Fprintln(tw, "-------\t----\t-----------\t----------\t----------")
	fmt.Fprintf(tw, "%s\t%.1f\t%.1f\t%d\t%s\n",
		e.Service,
		e.CPU,
		e.MemoryMB,
		e.UptimeSec,
		e.Timestamp.Format("15:04:05"),
	)
	tw.Flush()
}
