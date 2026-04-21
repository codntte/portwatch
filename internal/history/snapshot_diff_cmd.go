package history

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// PrintSnapshotDiffs writes a formatted table of snapshot diff entries to stdout.
func PrintSnapshotDiffs(path string, since time.Time, host string) error {
	return printSnapshotDiffsTo(os.Stdout, path, since, host)
}

func printSnapshotDiffsTo(w io.Writer, path string, since time.Time, host string) error {
	store := NewSnapshotDiffStore(path)
	entries, err := store.Load()
	if err != nil {
		return fmt.Errorf("snapshot_diff: load: %w", err)
	}
	if len(entries) == 0 {
		fmt.Fprintln(w, "No snapshot diff entries found.")
		return nil
	}
	fmt.Fprintf(w, "%-28s %-20s %-20s %s\n", "TIMESTAMP", "HOST", "OPENED", "CLOSED")
	fmt.Fprintln(w, strings.Repeat("-", 80))
	for _, e := range entries {
		if !since.IsZero() && e.Timestamp.Before(since) {
			continue
		}
		if host != "" && e.Host != host {
			continue
		}
		fmt.Fprintf(w, "%-28s %-20s %-20s %s\n",
			e.Timestamp.Format(time.RFC3339),
			e.Host,
			formatSnapshotPorts(e.Opened),
			formatSnapshotPorts(e.Closed),
		)
	}
	return nil
}

func formatSnapshotPorts(ports []int) string {
	if len(ports) == 0 {
		return "-"
	}
	parts := make([]string, len(ports))
	for i, p := range ports {
		parts[i] = fmt.Sprintf("%d", p)
	}
	return strings.Join(parts, ",")
}
