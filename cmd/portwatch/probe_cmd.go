package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/portwatch/internal/history"
)

func probeFile(dataDir string) string {
	return filepath.Join(dataDir, "probes.jsonl")
}

func runProbe(args []string) {
	fs := flag.NewFlagSet("probe", flag.ExitOnError)
	host := fs.String("host", "", "filter by host")
	sinceDur := fs.Duration("since", 0, "show probes newer than duration (e.g. 1h, 30m)")
	dataDir := fs.String("data", ".", "directory containing probe data")
	_ = fs.Parse(args)

	var since time.Time
	if *sinceDur > 0 {
		since = time.Now().UTC().Add(-*sinceDur)
	}

	path := probeFile(*dataDir)
	if err := history.PrintProbes(path, *host, since); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
