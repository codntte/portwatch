package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/portwatch/internal/history"
)

func snapshotDiffFile(dataDir string) string {
	return filepath.Join(dataDir, "snapshot_diff.jsonl")
}

func runSnapshotDiff(args []string) {
	fs := flag.NewFlagSet("snapshot-diff", flag.ExitOnError)
	dataDir := fs.String("data", ".", "directory containing snapshot diff data")
	host := fs.String("host", "", "filter by host address")
	sinceStr := fs.String("since", "", "show entries after this time (RFC3339)")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: portwatch snapshot-diff [flags]")
		fmt.Fprintln(os.Stderr, "\nFlags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "snapshot-diff: %v\n", err)
		os.Exit(1)
	}

	var since time.Time
	if *sinceStr != "" {
		var err error
		since, err = time.Parse(time.RFC3339, *sinceStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "snapshot-diff: invalid --since value %q: %v\n", *sinceStr, err)
			os.Exit(1)
		}
	}

	path := snapshotDiffFile(*dataDir)
	if err := history.PrintSnapshotDiffs(path, since, *host); err != nil {
		fmt.Fprintf(os.Stderr, "snapshot-diff: %v\n", err)
		os.Exit(1)
	}
}
