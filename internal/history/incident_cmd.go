package history

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

// PrintIncidents writes a formatted table of incidents to stdout.
func PrintIncidents(path string, showResolved bool) error {
	return printIncidentsTo(os.Stdout, path, showResolved)
}

func printIncidentsTo(w io.Writer, path string, showResolved bool) error {
	store := NewIncidentStore(path)
	incidents, err := store.Load()
	if err != nil {
		return fmt.Errorf("load incidents: %w", err)
	}

	var filtered []Incident
	for _, inc := range incidents {
		if !showResolved && inc.ResolvedAt != nil {
			continue
		}
		filtered = append(filtered, inc)
	}

	if len(filtered) == 0 {
		fmt.Fprintln(w, "no incidents found")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tHOST\tKIND\tSEVERITY\tPORTS\tCREATED\tRESOLVED")
	for _, inc := range filtered {
		resolved := "-"
		if inc.ResolvedAt != nil {
			resolved = inc.ResolvedAt.Format(time.RFC3339)
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			inc.ID,
			inc.Host,
			inc.Kind,
			string(inc.Severity),
			formatIncidentPorts(inc.Ports),
			inc.CreatedAt.Format(time.RFC3339),
			resolved,
		)
	}
	return tw.Flush()
}

func formatIncidentPorts(ports []int) string {
	if len(ports) == 0 {
		return "-"
	}
	parts := make([]string, len(ports))
	for i, p := range ports {
		parts[i] = fmt.Sprintf("%d", p)
	}
	return strings.Join(parts, ",")
}
