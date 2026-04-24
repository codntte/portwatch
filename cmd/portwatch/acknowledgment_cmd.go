package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/user/portwatch/internal/history"
)

func acknowledgmentFile(dir string) string {
	return filepath.Join(dir, "acknowledgments.json")
}

func runAcknowledgment(args []string) {
	fs := flag.NewFlagSet("ack", flag.ExitOnError)
	dataDir := fs.String("data", ".", "data directory")
	host := fs.String("host", "", "host to filter or target")
	port := fs.Int("port", 0, "port to acknowledge or delete")
	ackedBy := fs.String("by", "", "name of acknowledger")
	comment := fs.String("comment", "", "acknowledgment comment")
	ttl := fs.Duration("ttl", 0, "optional expiry duration (e.g. 24h)")
	delCmd := fs.Bool("delete", false, "delete acknowledgment for host:port")
	_ = fs.Parse(args)

	store := history.NewAcknowledgmentStore(acknowledgmentFile(*dataDir))

	if *delCmd {
		if *host == "" || *port == 0 {
			fmt.Fprintln(os.Stderr, "ack --delete requires --host and --port")
			os.Exit(1)
		}
		if err := store.Delete(*host, *port); err != nil {
			fmt.Fprintf(os.Stderr, "error deleting acknowledgment: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("acknowledgment removed for %s:%d\n", *host, *port)
		return
	}

	if fs.NArg() == 0 && *host == "" {
		// list mode
		entries, err := store.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading acknowledgments: %v\n", err)
			os.Exit(1)
		}
		if len(entries) == 0 {
			fmt.Println("no acknowledgments recorded")
			return
		}
		fmt.Printf("%-20s %-6s %-16s %-26s %s\n", "HOST", "PORT", "ACKED BY", "ACKED AT", "COMMENT")
		for _, a := range entries {
			fmt.Printf("%-20s %-6s %-16s %-26s %s\n",
				a.Host,
				strconv.Itoa(a.Port),
				a.AckedBy,
				a.AckedAt.Format(time.RFC3339),
				a.Comment,
			)
		}
		return
	}

	// add mode
	if *host == "" || *port == 0 || *ackedBy == "" {
		fmt.Fprintln(os.Stderr, "ack requires --host, --port, and --by")
		os.Exit(1)
	}
	a := history.Acknowledgment{
		ID:      fmt.Sprintf("%s-%d-%d", *host, *port, time.Now().UnixNano()),
		Host:    *host,
		Port:    *port,
		AckedBy: *ackedBy,
		Comment: *comment,
		AckedAt: time.Now().UTC(),
	}
	if *ttl > 0 {
		a.ExpiresAt = time.Now().UTC().Add(*ttl)
	}
	if err := store.Append(a); err != nil {
		fmt.Fprintf(os.Stderr, "error saving acknowledgment: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("acknowledged %s:%d by %s\n", *host, *port, *ackedBy)
}
