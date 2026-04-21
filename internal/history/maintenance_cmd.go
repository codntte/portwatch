package history

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// PrintMaintenance writes maintenance windows in a tabular format to stdout.
func PrintMaintenance(path string, host string) error {
	return printMaintenanceTo(os.Stdout, path, host)
}

func printMaintenanceTo(w io.Writer, path string, host string) error {
	store := NewMaintenanceStore(path)
	windows, err := store.Load()
	if err != nil {
		return err
	}

	var filtered []MaintenanceWindow
	for _, mw := range windows {
		if host == "" || mw.Host == host {
			filtered = append(filtered, mw)
		}
	}

	if len(filtered) == 0 {
		fmt.Fprintln(w, "no maintenance windows found")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tHOST\tSTARTS AT\tENDS AT\tACTIVE\tREASON")
	for _, mw := range filtered {
		active := "no"
		if mw.IsActive() {
			active = "yes"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n",
			mw.ID,
			mw.Host,
			mw.StartsAt.Format(time.RFC3339),
			mw.EndsAt.Format(time.RFC3339),
			active,
			mw.Reason,
		)
	}
	return tw.Flush()
}
