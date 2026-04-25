package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/portwatch/internal/history"
)

func dependencyFile(dir string) string {
	return dir + "/dependencies.json"
}

func runDependency(args []string, dataDir string) {
	fs := flag.NewFlagSet("dependency", flag.ExitOnError)
	action := fs.String("action", "list", "Action: list, add, delete")
	id := fs.String("id", "", "Dependency ID")
	host := fs.String("host", "", "Host that has the dependency")
	dependsOn := fs.String("depends-on", "", "Host being depended upon")
	ports := fs.String("ports", "", "Comma-separated port numbers (optional)")
	note := fs.String("note", "", "Optional note")
	fs.Parse(args)

	store := history.NewDependencyStore(dependencyFile(dataDir))

	switch *action {
	case "list":
		filterHost := fs.Arg(0)
		var deps []history.Dependency
		var err error
		if filterHost != "" {
			deps, err = store.LoadByHost(filterHost)
		} else {
			deps, err = store.Load()
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if len(deps) == 0 {
			fmt.Println("no dependencies recorded")
			return
		}
		fmt.Printf("%-12s %-20s %-20s %-16s %s\n", "ID", "HOST", "DEPENDS ON", "PORTS", "NOTE")
		for _, d := range deps {
			portStr := formatDepPorts(d.Ports)
			fmt.Printf("%-12s %-20s %-20s %-16s %s\n", d.ID, d.Host, d.DependsOn, portStr, d.Note)
		}

	case "add":
		if *id == "" || *host == "" || *dependsOn == "" {
			fmt.Fprintln(os.Stderr, "error: --id, --host, and --depends-on are required")
			os.Exit(1)
		}
		dep := history.Dependency{
			ID:        *id,
			Host:      *host,
			DependsOn: *dependsOn,
			Ports:     parseDepPorts(*ports),
			Note:      *note,
		}
		if err := store.Add(dep); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("dependency %s added\n", *id)

	case "delete":
		if *id == "" {
			fmt.Fprintln(os.Stderr, "error: --id is required")
			os.Exit(1)
		}
		if err := store.Delete(*id); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("dependency %s deleted\n", *id)

	default:
		fmt.Fprintf(os.Stderr, "unknown action: %s\n", *action)
		os.Exit(1)
	}
}

func parseDepPorts(s string) []int {
	if s == "" {
		return nil
	}
	var ports []int
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if n, err := strconv.Atoi(p); err == nil {
			ports = append(ports, n)
		}
	}
	return ports
}

func formatDepPorts(ports []int) string {
	if len(ports) == 0 {
		return "-"
	}
	parts := make([]string, len(ports))
	for i, p := range ports {
		parts[i] = strconv.Itoa(p)
	}
	return strings.Join(parts, ",")
}
