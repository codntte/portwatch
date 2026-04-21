package history

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// PrintPolicies writes a formatted table of policies to stdout.
func PrintPolicies(path string) error {
	return printPoliciesTo(os.Stdout, path)
}

func printPoliciesTo(w io.Writer, path string) error {
	store := NewPolicyStore(path)
	entries, err := store.Load()
	if err != nil {
		return fmt.Errorf("policy: %w", err)
	}
	if len(entries) == 0 {
		fmt.Fprintln(w, "no policies defined")
		return nil
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NAME\tHOST\tPORT RANGE\tINTERVAL\tENABLED\tCREATED")
	for _, e := range entries {
		enabled := "no"
		if e.Enabled {
			enabled = "yes"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n",
			e.Name,
			e.Host,
			e.PortRange,
			formatPolicyInterval(e.Interval),
			enabled,
			e.CreatedAt.Format(time.RFC3339),
		)
	}
	return tw.Flush()
}

func formatPolicyInterval(d time.Duration) string {
	if d == 0 {
		return "—"
	}
	return d.String()
}
