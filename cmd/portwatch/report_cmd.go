package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/portwatch/internal/history"
)

func reportDiffFile(dataDir string) string {
	return filepath.Join(dataDir, "diffs.jsonl")
}

func reportStatsFile(dataDir string) string {
	return filepath.Join(dataDir, "stats.jsonl")
}

func runReport(args []string) {
	dataDir := os.Getenv("PORTWATCH_DATA_DIR")
	if dataDir == "" {
		dataDir = "."
	}

	saveJSON := false
	outPath := ""

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
			saveJSON = true
		case "--out":
			if i+1 < len(args) {
				i++
				outPath = args[i]
			}
		default:
			fmt.Fprintf(os.Stderr, "unknown flag: %s\n", args[i])
			os.Exit(1)
		}
	}

	diffFile := reportDiffFile(dataDir)
	statsFile := reportStatsFile(dataDir)

	if saveJSON {
		if outPath == "" {
			outPath = filepath.Join(dataDir, "report.json")
		}
		report, err := history.BuildReport(diffFile, statsFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error building report: %v\n", err)
			os.Exit(1)
		}
		if err := history.SaveReport(outPath, report); err != nil {
			fmt.Fprintf(os.Stderr, "error saving report: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Report saved to %s\n", outPath)
		return
	}

	if err := history.PrintReport(diffFile, statsFile); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
