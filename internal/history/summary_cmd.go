package history

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// PrintSummary writes a formatted summary table to w.
func PrintSummary(w io.Writer, summaries []HostSummary) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "HOST\tEVENTS\tOPENED\tCLOSED\tLAST SEEN")
	fmt.Fprintln(tw, "----\t------\t------\t------\t---------")
	for _, s := range summaries {
		lastSeen := "-"
		if !s.LastSeen.IsZero() {
			lastSeen = s.LastSeen.Format(time.RFC3339)
		}
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%s\n",
			s.Host, s.TotalEvents, s.Opened, s.Closed, lastSeen)
	}
	return tw.Flush()
}
