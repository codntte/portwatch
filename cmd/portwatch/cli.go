package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/user/portwatch/internal/history"
)

// usage prints help text to stdout.
func usage() {
	fmt.Println(`portwatch - monitor and alert on port changes

Usage:
  portwatch [flags]
  portwatch <command> [args]

Commands:
  history     Show scan history
  summary     Show port change summary
  export      Export history to CSV

Flags:
  -config <path>   Path to config file (default: portwatch.yaml)
  -once            Run a single scan and exit
  -help            Show this help message

Examples:
  portwatch
  portwatch -config /etc/portwatch.yaml
  portwatch -once
  portwatch history -host 192.168.1.1 -limit 20
  portwatch summary
  portwatch export -out history.csv`)
}

// runCLI parses os.Args and dispatches sub-commands.
// Returns an exit code.
func runCLI(args []string, store *history.Store) int {
	if len(args) == 0 {
		return 0 // caller handles default run
	}

	switch args[0] {
	case "help", "-help", "--help":
		usage()
		return 0

	case "summary":
		if err := history.PrintSummary(store); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		return 0

	case "export":
		out := "history.csv"
		for i, a := range args[1:] {
			if a == "-out" && i+1 < len(args)-1 {
				out = args[i+2]
			}
		}
		if err := store.ExportCSV(out); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		fmt.Printf("exported history to %s\n", out)
		return 0

	case "history":
		var host string
		var limit int
		for i, a := range args[1:] {
			if a == "-host" && i+1 < len(args)-1 {
				host = args[i+2]
			}
			if a == "-limit" && i+1 < len(args)-1 {
				n, err := strconv.Atoi(args[i+2])
				if err == nil {
					limit = n
				}
			}
		}
		entries, err := store.Query(history.QueryOptions{
			Host:  host,
			Limit: limit,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		for _, e := range entries {
			fmt.Printf("%s  %-20s  opened=%v closed=%v\n",
				e.Timestamp.Format("2006-01-02 15:04:05"),
				e.Host,
				e.Opened,
				e.Closed,
			)
		}
		return 0
	}

	fmt.Fprintf(os.Stderr, "unknown command: %s\nRun 'portwatch help' for usage.\n", args[0])
	return 1
}
