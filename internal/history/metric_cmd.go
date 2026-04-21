package history

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// PrintMetrics writes metric entries to stdout in a tabular format.
func PrintMetrics(path, host, name string) error {
	return printMetricsTo(os.Stdout, path, host, name)
}

func printMetricsTo(w io.Writer, path, host, name string) error {
	store := NewMetricStore(path)

	var entries []MetricEntry
	var err error
	if host != "" {
		entries, err = store.LoadByHost(host, name)
	} else {
		entries, err = store.Load()
	}
	if err != nil {
		return fmt.Errorf("load metrics: %w", err)
	}

	if len(entries) == 0 {
		fmt.Fprintln(w, "no metric entries found")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tHOST\tNAME\tVALUE\tUNIT")
	for _, e := range entries {
		unit := e.Unit
		if unit == "" {
			unit = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%.4g\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Host,
			e.Name,
			e.Value,
			unit,
		)
	}
	return tw.Flush()
}
