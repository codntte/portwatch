package main

import (
	"fmt"
	"os"

	"github.com/user/portwatch/internal/history"
)

// runStats prints port open/close frequency stats for a given host.
func runStats(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: portwatch stats <host>")
		os.Exit(1)
	}
	host := args[0]

	file := historyFile()
	store := history.NewStore(file)

	stats, err := store.Stats(host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading history: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Port statistics for %s:\n", host)
	fmt.Print(history.FormatStats(stats))
}
