package main

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/user/portwatch/internal/history"
)

func watchlistFile(dataDir string) string {
	return dataDir + "/watchlist.json"
}

// runWatchlist handles the `watchlist` subcommand.
// Usage:
//
//	portwatch watchlist add <host> <port> [label]
//	portwatch watchlist remove <host> <port>
//	portwatch watchlist list
func runWatchlist(args []string, dataDir string) {
	store := history.NewWatchlistStore(watchlistFile(dataDir))

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: portwatch watchlist <add|remove|list> [args]")
		os.Exit(1)
	}

	switch args[0] {
	case "add":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: portwatch watchlist add <host> <port> [label]")
			os.Exit(1)
		}
		port, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid port %q: %v\n", args[2], err)
			os.Exit(1)
		}
		label := ""
		if len(args) >= 4 {
			label = args[3]
		}
		if err := store.Add(args[1], port, label); err != nil {
			fmt.Fprintf(os.Stderr, "error adding to watchlist: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("added %s:%d to watchlist\n", args[1], port)

	case "remove":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: portwatch watchlist remove <host> <port>")
			os.Exit(1)
		}
		port, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid port %q: %v\n", args[2], err)
			os.Exit(1)
		}
		if err := store.Delete(args[1], port); err != nil {
			fmt.Fprintf(os.Stderr, "error removing from watchlist: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("removed %s:%d from watchlist\n", args[1], port)

	case "list":
		entries, err := store.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading watchlist: %v\n", err)
			os.Exit(1)
		}
		if len(entries) == 0 {
			fmt.Println("watchlist is empty")
			return
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "HOST\tPORT\tLABEL\tADDED")
		for _, e := range entries {
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\n",
				e.Host, e.Port, e.Label, e.AddedAt.Format("2006-01-02 15:04:05"))
		}
		w.Flush()

	default:
		fmt.Fprintf(os.Stderr, "unknown watchlist command: %q\n", args[0])
		os.Exit(1)
	}
}
