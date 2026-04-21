package history

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

// PrintAlertLog writes a formatted alert log to stdout.
func PrintAlertLog(path string, since time.Time, host string) error {
	return printAlertLogTo(os.Stdout, path, since, host)
}

func printAlertLogTo(w io.Writer, path string, since time.Time, host string) error {
	store := NewAlertLogStore(path)
	entries, err := store.Load()
	if err != nil {
		return fmt.Errorf("loading alert log: %w", err)
	}

	var filtered []AlertEntry
	for _, e := range entries {
		if !since.IsZero() && e.Timestamp.Before(since) {
			continue
		}
		if host != "" && e.Host != host {
			continue
		}
		filtered = append(filtered, e)
	}

	if len(filtered) == 0 {
		fmt.Fprintln(w, "no alert entries found")
		return nil
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.Before(filtered[j].Timestamp)
	})

	fmt.Fprintf(w, "%-25s %-20s %-20s %s\n", "TIMESTAMP", "HOST", "OPENED", "CLOSED")
	fmt.Fprintln(w, strings.Repeat("-", 80))
	for _, e := range filtered {
		ts := e.Timestamp.Format(time.RFC3339)
		opened := formatAlertPorts(e.Opened)
		closed := formatAlertPorts(e.Closed)
		fmt.Fprintf(w, "%-25s %-20s %-20s %s\n", ts, e.Host, opened, closed)
	}
	return nil
}

func formatAlertPorts(ports []int) string {
	if len(ports) == 0 {
		return "-"
	}
	parts := make([]string, len(ports))
	for i, p := range ports {
		parts[i] = fmt.Sprintf("%d", p)
	}
	return strings.Join(parts, ",")
}
