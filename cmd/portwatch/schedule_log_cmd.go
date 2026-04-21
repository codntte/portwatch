package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/portwatch/internal/history"
)

func scheduleLogFile(dataDir string) string {
	return filepath.Join(dataDir, "schedule_log.jsonl")
}

func runScheduleLog(args []string, dataDir string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: portwatch schedule-log <list>")
		os.Exit(1)
	}

	switch args[0] {
	case "list":
		if err := history.PrintScheduleLog(scheduleLogFile(dataDir)); err != nil {
			fmt.Fprintf(os.Stderr, "schedule-log list: %v\n", err)
			os.Exit(1)
		}
	case "clear":
		path := scheduleLogFile(dataDir)
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "schedule-log clear: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Schedule log cleared.")
	default:
		fmt.Fprintf(os.Stderr, "unknown schedule-log command: %s\n", args[0])
		os.Exit(1)
	}
}
