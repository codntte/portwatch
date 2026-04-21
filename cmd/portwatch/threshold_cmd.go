package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/user/portwatch/internal/history"
)

func thresholdFile(dataDir string) string {
	return filepath.Join(dataDir, "thresholds.json")
}

func runThreshold(args []string, dataDir string) {
	path := thresholdFile(dataDir)
	store := history.NewThresholdStore(path)

	if len(args) == 0 {
		if err := history.PrintThresholds(path); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	switch args[0] {
	case "add":
		if len(args) < 5 {
			fmt.Fprintln(os.Stderr, "usage: threshold add <host> <port> <max_closed> <window>")
			os.Exit(1)
		}
		port, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid port: %v\n", err)
			os.Exit(1)
		}
		maxClosed, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid max_closed: %v\n", err)
			os.Exit(1)
		}
		rule := history.ThresholdRule{
			Host:      args[1],
			Port:      port,
			MaxClosed: maxClosed,
			Window:    args[4],
		}
		if err := store.Add(rule); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("threshold added for %s:%d\n", rule.Host, rule.Port)

	case "delete":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: threshold delete <host> <port>")
			os.Exit(1)
		}
		port, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid port: %v\n", err)
			os.Exit(1)
		}
		if err := store.Delete(args[1], port); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("threshold deleted for %s:%d\n", args[1], port)

	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", args[0])
		os.Exit(1)
	}
}
