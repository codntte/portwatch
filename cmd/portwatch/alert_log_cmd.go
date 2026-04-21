package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/portwatch/internal/history"
)

// alertLogFile returns the path to the alert log file within the data directory.
func alertLogFile(dataDir string) string {
	return filepath.Join(dataDir, "alert_log.json")
}

// runAlertLog handles the `alert-log` subcommand, which displays the history
// of alert events recorded during port scan cycles.
//
// Usage:
//
//	portwatch alert-log [flags]
//
// Flags:
//
//	-data   path to the data directory (default: .portwatch)
//	-host   filter entries by host address
//	-limit  maximum number of entries to display (0 = all)
func runAlertLog(args []string) {
	fs := flag.NewFlagSet("alert-log", flag.ExitOnError)

	dataDir := fs.String("data", ".portwatch", "path to data directory")
	host := fs.String("host", "", "filter by host address")
	limit := fs.Int("limit", 0, "maximum number of entries to show (0 = all)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "alert-log: failed to parse flags: %v\n", err)
		os.Exit(1)
	}

	path := alertLogFile(*dataDir)
	store := history.NewAlertLogStore(path)

	entries, err := store.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "alert-log: failed to load entries: %v\n", err)
		os.Exit(1)
	}

	// Apply host filter.
	if *host != "" {
		filtered := entries[:0]
		for _, e := range entries {
			if e.Host == *host {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	// Apply limit (take from the end to show most recent).
	if *limit > 0 && len(entries) > *limit {
		entries = entries[len(entries)-*limit:]
	}

	history.PrintAlertLog(os.Stdout, entries)
}
