package history

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

// PrintEvents prints all WatchEvents from the store to stdout in a table.
func PrintEvents(storePath string, since time.Time) error {
	s := &Store{path: storePath}
	events, err := s.LoadEvents()
	if err != nil {
		return fmt.Errorf("load events: %w", err)
	}
	if len(events) == 0 {
		fmt.Println("no events recorded")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tHOST\tOPENED\tCLOSED")

	for _, e := range events {
		if !since.IsZero() && e.Timestamp.Before(since) {
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%v\t%v\n",
			e.Timestamp.Format(time.RFC3339),
			e.Host,
			formatPorts(e.Opened),
			formatPorts(e.Closed),
		)
	}
	return w.Flush()
}

func formatPorts(ports []int) string {
	if len(ports) == 0 {
		return "-"
	}
	s := ""
	for i, p := range ports {
		if i > 0 {
			s += ","
		}
		s += fmt.Sprintf("%d", p)
	}
	return s
}
