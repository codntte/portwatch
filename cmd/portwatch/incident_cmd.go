package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/portwatch/internal/history"
)

func incidentFile(dataDir string) string {
	return filepath.Join(dataDir, "incidents.json")
}

func runIncident(args []string) {
	fs := flag.NewFlagSet("incident", flag.ExitOnError)
	var (
		list        = fs.Bool("list", false, "list open incidents")
		listAll     = fs.Bool("list-all", false, "list all incidents including resolved")
		resolveID   = fs.String("resolve", "", "mark incident as resolved by ID")
		addHost     = fs.String("host", "", "host for new incident")
		addKind     = fs.String("kind", "opened", "kind: opened or closed")
		addSeverity = fs.String("severity", "medium", "severity: low, medium, high")
		addNote     = fs.String("note", "", "optional note")
		addPorts    = fs.String("ports", "", "comma-separated port list")
		dataDir     = fs.String("data-dir", ".", "directory for data files")
	)
	_ = fs.Parse(args)

	path := incidentFile(*dataDir)

	switch {
	case *list || *listAll:
		if err := history.PrintIncidents(path, *listAll); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case *resolveID != "":
		store := history.NewIncidentStore(path)
		if err := store.Resolve(*resolveID); err != nil {
			fmt.Fprintf(os.Stderr, "error resolving incident: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("incident %s marked as resolved\n", *resolveID)

	case *addHost != "":
		ports := parseIncidentPorts(*addPorts)
		inc := history.Incident{
			ID:        fmt.Sprintf("inc-%d", time.Now().UnixNano()),
			Host:      *addHost,
			Kind:      *addKind,
			Severity:  history.IncidentSeverity(*addSeverity),
			Note:      *addNote,
			Ports:     ports,
			CreatedAt: time.Now().UTC(),
		}
		store := history.NewIncidentStore(path)
		if err := store.Append(inc); err != nil {
			fmt.Fprintf(os.Stderr, "error appending incident: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("incident %s recorded\n", inc.ID)

	default:
		fs.Usage()
		os.Exit(1)
	}
}

func parseIncidentPorts(s string) []int {
	if s == "" {
		return nil
	}
	var ports []int
	for _, part := range splitCSV(s) {
		var p int
		if _, err := fmt.Sscanf(part, "%d", &p); err == nil {
			ports = append(ports, p)
		}
	}
	return ports
}

func splitCSV(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			if tok := s[start:i]; tok != "" {
				out = append(out, tok)
			}
			start = i + 1
		}
	}
	return out
}
