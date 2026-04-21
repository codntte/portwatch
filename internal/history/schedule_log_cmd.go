package history

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// PrintScheduleLog writes a formatted table of schedule log entries to stdout.
func PrintScheduleLog(path string) error {
	return printScheduleLogTo(os.Stdout, path)
}

func printScheduleLogTo(w io.Writer, path string) error {
	store := NewScheduleLogStore(path)
	entries, err := store.Load()
	if err != nil {
		return fmt.Errorf("loading schedule log: %w", err)
	}

	if len(entries) == 0 {
		fmt.Fprintln(w, "No schedule log entries found.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tHOST\tDURATION\tPORTS OPEN\tERROR")
	for _, e := range entries {
		errStr := e.Error
		if errStr == "" {
			errStr = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Host,
			e.Duration,
			e.PortsOpen,
			errStr,
		)
	}
	return tw.Flush()
}
