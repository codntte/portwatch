package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/portwatch/internal/history"
)

func baselineFile(dataDir string) string {
	return dataDir + "/baselines.json"
}

func runBaseline(args []string, dataDir string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: portwatch baseline <add|get|delete|list> [args...]")
		os.Exit(1)
	}

	store := history.NewBaselineStore(baselineFile(dataDir))
	subcmd := args[0]

	switch subcmd {
	case "add":
		if len(args) < 4 {
			fmt.Fprintln(os.Stderr, "usage: portwatch baseline add <name> <host> <port,...>")
			os.Exit(1)
		}
		name, host, portStr := args[1], args[2], args[3]
		ports, err := parsePorts(portStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid ports %q: %v\n", portStr, err)
			os.Exit(1)
		}
		if err := store.Add(name, host, ports); err != nil {
			fmt.Fprintf(os.Stderr, "error saving baseline: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("baseline %q saved for host %s\n", name, host)

	case "get":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: portwatch baseline get <name> <host>")
			os.Exit(1)
		}
		b, ok, err := store.Get(args[1], args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if !ok {
			fmt.Printf("no baseline %q found for host %s\n", args[1], args[2])
			return
		}
		fmt.Printf("name=%s host=%s ports=%v created=%s\n",
			b.Name, b.Host, b.Ports, b.CreatedAt.Format("2006-01-02 15:04:05"))

	case "delete":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: portwatch baseline delete <name> <host>")
			os.Exit(1)
		}
		if err := store.Delete(args[1], args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("baseline %q deleted for host %s\n", args[1], args[2])

	case "list":
		list, err := store.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if len(list) == 0 {
			fmt.Println("no baselines stored")
			return
		}
		for _, b := range list {
			fmt.Printf("%-20s %-20s ports=%-30v created=%s\n",
				b.Name, b.Host, b.Ports, b.CreatedAt.Format("2006-01-02 15:04:05"))
		}

	default:
		fmt.Fprintf(os.Stderr, "unknown baseline subcommand: %s\n", subcmd)
		os.Exit(1)
	}
}

func parsePorts(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	ports := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid port %q", p)
		}
		ports = append(ports, n)
	}
	return ports, nil
}
