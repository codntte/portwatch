package history

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// PrintThresholds writes all threshold rules to stdout in a table format.
func PrintThresholds(path string) error {
	return printThresholdsTo(os.Stdout, path)
}

func printThresholdsTo(w io.Writer, path string) error {
	store := NewThresholdStore(path)
	rules, err := store.Load()
	if err != nil {
		return fmt.Errorf("threshold: load: %w", err)
	}
	if len(rules) == 0 {
		fmt.Fprintln(w, "no threshold rules defined")
		return nil
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "HOST\tPORT\tMAX_CLOSED\tWINDOW\tCREATED")
	for _, r := range rules {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%s\t%s\n",
			r.Host,
			r.Port,
			r.MaxClosed,
			r.Window,
			r.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	}
	return tw.Flush()
}
