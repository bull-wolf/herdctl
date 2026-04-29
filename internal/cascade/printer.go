package cascade

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Print writes a human-readable table of cascade rules to w.
func Print(m *Manager, w io.Writer) {
	rules := m.All()
	if len(rules) == 0 {
		fmt.Fprintln(w, "no cascade rules registered")
		return
	}
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Service < rules[j].Service
	})
	fmt.Fprintf(w, "%-20s %-10s %s\n", "SERVICE", "POLICY", "TARGETS")
	fmt.Fprintln(w, strings.Repeat("-", 55))
	for _, r := range rules {
		targets := strings.Join(r.Targets, ", ")
		if targets == "" {
			targets = "-"
		}
		fmt.Fprintf(w, "%-20s %-10s %s\n", r.Service, string(r.Policy), targets)
	}
}
