package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/user/portwatch/internal/history"
)

func silenceFile(dataDir string) string {
	return filepath.Join(dataDir, "silenced.json")
}

func runSilence(args []string, dataDir string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: portwatch silence <add|delete|list> [host] [port] [reason] [--ttl=<duration>]")
		os.Exit(1)
	}

	store := history.NewSilencedStore(silenceFile(dataDir))

	switch args[0] {
	case "list":
		entries, err := store.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if len(entries) == 0 {
			fmt.Println("no active silence rules")
			return
		}
		fmt.Printf("%-20s  %-6s  %-20s  %s\n", "HOST", "PORT", "EXPIRES", "REASON")
		for _, e := range entries {
			expires := "never"
			if !e.ExpiresAt.IsZero() {
				expires = e.ExpiresAt.Format(time.RFC3339)
			}
			fmt.Printf("%-20s  %-6d  %-20s  %s\n", e.Host, e.Port, expires, e.Reason)
		}

	case "add":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: portwatch silence add <host> <port> [reason] [--ttl=<duration>]")
			os.Exit(1)
		}
		host := args[1]
		port, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid port: %v\n", err)
			os.Exit(1)
		}
		reason := ""
		if len(args) >= 4 {
			reason = args[3]
		}
		var expiresAt time.Time
		for _, a := range args[3:] {
			var ttlStr string
			if _, err := fmt.Sscanf(a, "--ttl=%s", &ttlStr); err == nil {
				d, err := time.ParseDuration(ttlStr)
				if err != nil {
					fmt.Fprintf(os.Stderr, "invalid ttl: %v\n", err)
					os.Exit(1)
				}
				expiresAt = time.Now().Add(d)
			}
		}
		entry := history.SilencedEntry{
			Host:      host,
			Port:      port,
			Reason:    reason,
			CreatedAt: time.Now(),
			ExpiresAt: expiresAt,
		}
		if err := store.Add(entry); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("silenced %s:%d\n", host, port)

	case "delete":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: portwatch silence delete <host> <port>")
			os.Exit(1)
		}
		host := args[1]
		port, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid port: %v\n", err)
			os.Exit(1)
		}
		if err := store.Delete(host, port); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("removed silence rule for %s:%d\n", host, port)

	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", args[0])
		os.Exit(1)
	}
}
