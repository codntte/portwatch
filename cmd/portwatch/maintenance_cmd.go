package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/portwatch/internal/history"
)

func maintenanceFile(dir string) string {
	return filepath.Join(dir, "maintenance.jsonl")
}

func runMaintenance(args []string, dataDir string) error {
	fs := flag.NewFlagSet("maintenance", flag.ContinueOnError)
	add := fs.Bool("add", false, "add a new maintenance window")
	del := fs.String("delete", "", "delete window by ID")
	list := fs.Bool("list", false, "list maintenance windows")
	host := fs.String("host", "", "filter by host")
	id := fs.String("id", "", "window ID (required with --add)")
	reason := fs.String("reason", "", "reason for maintenance")
	start := fs.String("start", "", "start time (RFC3339)")
	end := fs.String("end", "", "end time (RFC3339)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	path := maintenanceFile(dataDir)
	store := history.NewMaintenanceStore(path)

	switch {
	case *add:
		if *id == "" || *start == "" || *end == "" || *host == "" {
			return fmt.Errorf("--add requires --id, --host, --start, and --end")
		}
		startsAt, err := time.Parse(time.RFC3339, *start)
		if err != nil {
			return fmt.Errorf("invalid --start: %w", err)
		}
		endsAt, err := time.Parse(time.RFC3339, *end)
		if err != nil {
			return fmt.Errorf("invalid --end: %w", err)
		}
		w := history.MaintenanceWindow{
			ID:        *id,
			Host:      *host,
			StartsAt:  startsAt,
			EndsAt:    endsAt,
			Reason:    *reason,
			CreatedAt: time.Now().UTC(),
		}
		if err := store.Add(w); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "maintenance window %s added\n", *id)

	case *del != "":
		if err := store.Delete(*del); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "maintenance window %s deleted\n", *del)

	case *list:
		return history.PrintMaintenance(path, *host)

	default:
		return history.PrintMaintenance(path, *host)
	}

	return nil
}
