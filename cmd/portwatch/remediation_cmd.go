package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/portwatch/internal/history"
)

func remediationFile(dir string) string {
	return filepath.Join(dir, "remediation.json")
}

func runRemediation(args []string, dataDir string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: portwatch remediation <list|add|update|delete> [flags]")
		os.Exit(1)
	}

	store := history.NewRemediationStore(remediationFile(dataDir))

	switch args[0] {
	case "list":
		host := ""
		if len(args) >= 3 && args[1] == "--host" {
			host = args[2]
		}
		if err := history.PrintRemediations(remediationFile(dataDir), host); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case "add":
		if len(args) < 5 {
			fmt.Fprintln(os.Stderr, "usage: portwatch remediation add <id> <host> <port> <action>")
			os.Exit(1)
		}
		var port int
		if _, err := fmt.Sscanf(args[3], "%d", &port); err != nil {
			fmt.Fprintf(os.Stderr, "invalid port: %s\n", args[3])
			os.Exit(1)
		}
		now := time.Now().UTC()
		e := history.RemediationEntry{
			ID:        args[1],
			Host:      args[2],
			Port:      port,
			Action:    args[4],
			Status:    history.RemediationPending,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := store.Append(e); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("remediation %s added\n", e.ID)

	case "update":
		if len(args) < 4 {
			fmt.Fprintln(os.Stderr, "usage: portwatch remediation update <id> <status> [note]")
			os.Exit(1)
		}
		note := ""
		if len(args) >= 5 {
			note = args[4]
		}
		if err := store.UpdateStatus(args[1], history.RemediationStatus(args[2]), note); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("remediation %s updated to %s\n", args[1], args[2])

	case "delete":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "usage: portwatch remediation delete <id>")
			os.Exit(1)
		}
		if err := store.Delete(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("remediation %s deleted\n", args[1])

	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", args[0])
		os.Exit(1)
	}
}
