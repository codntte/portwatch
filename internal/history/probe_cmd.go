package history

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// PrintProbes writes a formatted probe result table to stdout.
func PrintProbes(path, host string, since time.Time) error {
	return printProbesTo(os.Stdout, path, host, since)
}

func printProbesTo(w io.Writer, path, host string, since time.Time) error {
	store := NewProbeStore(path)

	var results []ProbeResult
	var err error
	if host != "" {
		results, err = store.LoadByHost(host)
	} else {
		results, err = store.Load()
	}
	if err != nil {
		return fmt.Errorf("load probes: %w", err)
	}

	if !since.IsZero() {
		var filtered []ProbeResult
		for _, r := range results {
			if !r.Timestamp.Before(since) {
				filtered = append(filtered, r)
			}
		}
		results = filtered
	}

	if len(results) == 0 {
		fmt.Fprintln(w, "no probe results found")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "HOST\tTIMESTAMP\tLATENCY (ms)\tSTATUS\tERROR")
	for _, r := range results {
		status := "ok"
		if !r.Success {
			status = "fail"
		}
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\t%s\n",
			r.Host,
			r.Timestamp.Format(time.RFC3339),
			r.LatencyMs,
			status,
			r.Error,
		)
	}
	return tw.Flush()
}
