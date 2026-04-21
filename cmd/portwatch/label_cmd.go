package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/user/portwatch/internal/history"
)

func labelFile(dataDir string) string {
	return filepath.Join(dataDir, "labels.json")
}

func runLabel(args []string, dataDir string) {
	store := history.NewLabelStore(labelFile(dataDir))

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: portwatch label <add|delete|list> [host] [name]")
		os.Exit(1)
	}

	switch args[0] {
	case "add":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: portwatch label add <host> <name>")
			os.Exit(1)
		}
		if err := store.Add(args[1], args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("label %q added to host %q\n", args[2], args[1])

	case "delete":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: portwatch label delete <host> <name>")
			os.Exit(1)
		}
		if err := store.Delete(args[1], args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("label %q removed from host %q\n", args[2], args[1])

	case "list":
		host := ""
		if len(args) >= 2 {
			host = args[1]
		}
		labels, err := store.Load(host)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		printLabels(labels)

	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

func printLabels(labels []history.Label) {
	if len(labels) == 0 {
		fmt.Println("no labels found")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "HOST\tLABEL\tCREATED")
	for _, l := range labels {
		fmt.Fprintf(w, "%s\t%s\t%s\n", l.Host, l.Name, l.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	w.Flush()
}
