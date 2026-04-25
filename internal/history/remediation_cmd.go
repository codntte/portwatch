package history

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// PrintRemediations writes remediation entries to stdout in a tabular format.
func PrintRemediations(path, filterHost string) error {
	return printRemediationsTo(os.Stdout, path, filterHost)
}

func printRemediationsTo(w io.Writer, path, filterHost string) error {
	store := NewRemediationStore(path)
	entries, err := store.Load()
	if err != nil {
		return fmt.Errorf("load remediations: %w", err)
	}

	if len(entries) == 0 {
		fmt.Fprintln(w, "no remediation entries found")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tHOST\tPORT\tACTION\tSTATUS\tNOTE\tUPDATED")

	for _, e := range entries {
		if filterHost != "" && e.Host != filterHost {
			continue
		}
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\t%s\t%s\t%s\n",
			e.ID,
			e.Host,
			e.Port,
			e.Action,
			e.Status,
			e.Note,
			e.UpdatedAt.Format("2006-01-02 15:04:05"),
		)
	}
	return tw.Flush()
}
