package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/user/portwatch/internal/history"
)

func bookmarkFile(dataDir string) string {
	return filepath.Join(dataDir, "bookmarks.json")
}

// runBookmark handles the `portwatch bookmark` subcommand.
//
// Usage:
//
//	portwatch bookmark add <name> <host> <ports> [note]
//	portwatch bookmark list [host]
//	portwatch bookmark delete <name> <host>
func runBookmark(args []string, dataDir string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: portwatch bookmark <add|list|delete> ...")
	}

	store := bookmarkFile(dataDir)
	s := history.NewBookmarkStore(store)

	switch args[0] {
	case "add":
		if len(args) < 4 {
			return fmt.Errorf("usage: portwatch bookmark add <name> <host> <ports> [note]")
		}
		name := args[1]
		host := args[2]
		ports, err := parseBookmarkPorts(args[3])
		if err != nil {
			return fmt.Errorf("invalid ports %q: %w", args[3], err)
		}
		note := ""
		if len(args) >= 5 {
			note = strings.Join(args[4:], " ")
		}
		b := history.Bookmark{
			Name:      name,
			Host:      host,
			CreatedAt: time.Now().UTC(),
			Note:      note,
			Ports:     ports,
		}
		if err := s.Add(b); err != nil {
			return fmt.Errorf("add bookmark: %w", err)
		}
		fmt.Fprintf(os.Stdout, "bookmark %q saved for host %s\n", name, host)

	case "list":
		filterHost := ""
		if len(args) >= 2 {
			filterHost = args[1]
		}
		return history.PrintBookmarks(store, filterHost)

	case "delete":
		if len(args) < 3 {
			return fmt.Errorf("usage: portwatch bookmark delete <name> <host>")
		}
		if err := s.Delete(args[1], args[2]); err != nil {
			return fmt.Errorf("delete bookmark: %w", err)
		}
		fmt.Fprintf(os.Stdout, "bookmark %q deleted\n", args[1])

	default:
		return fmt.Errorf("unknown bookmark command %q", args[0])
	}
	return nil
}

func parseBookmarkPorts(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	ports := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		ports = append(ports, n)
	}
	return ports, nil
}
