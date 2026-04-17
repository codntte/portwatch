package main

import (
	"fmt"
	"os"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/history"
)

// runRetention loads config, applies the retention policy to each host's
// history file, and prints a summary of what was done.
func runRetention(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	policy := history.RetentionPolicy{
		MaxAge:     cfg.Retention.MaxAge(),
		MaxEntries: cfg.Retention.MaxEntries,
	}

	for _, host := range cfg.Hosts {
		path := historyFile(host.Name)
		if err := policy.Apply(path); err != nil {
			fmt.Fprintf(os.Stderr, "retention: host %s: %v\n", host.Name, err)
			continue
		}
		fmt.Printf("retention applied: %s -> %s\n", host.Name, path)
	}
	return nil
}

// historyFile returns the path for a host's history file.
func historyFile(hostName string) string {
	return fmt.Sprintf(".portwatch/%s.history.jsonl", hostName)
}
