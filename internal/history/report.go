package history

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// ReportEntry summarizes scan activity for a single host.
type ReportEntry struct {
	Host      string    `json:"host"`
	OpenPorts []int     `json:"open_ports"`
	Changes   int       `json:"changes"`
	LastSeen  time.Time `json:"last_seen"`
}

// Report holds a full summary report across all hosts.
type Report struct {
	GeneratedAt time.Time     `json:"generated_at"`
	Entries     []ReportEntry `json:"entries"`
}

// BuildReport reads diffs and stats to produce a Report.
func BuildReport(diffFile, statsFile string) (*Report, error) {
	diffs, err := LoadDiffs(diffFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load diffs: %w", err)
	}

	changeCount := make(map[string]int)
	lastSeen := make(map[string]time.Time)
	for _, d := range diffs {
		changeCount[d.Host]++
		if d.Timestamp.After(lastSeen[d.Host]) {
			lastSeen[d.Host] = d.Timestamp
		}
	}

	stats, err := LoadStats(statsFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load stats: %w", err)
	}

	hostPorts := make(map[string][]int)
	for _, s := range stats {
		hostPorts[s.Host] = s.OpenPorts
	}

	hosts := make([]string, 0)
	seen := make(map[string]bool)
	for h := range changeCount {
		if !seen[h] {
			hosts = append(hosts, h)
			seen[h] = true
		}
	}
	for h := range hostPorts {
		if !seen[h] {
			hosts = append(hosts, h)
			seen[h] = true
		}
	}
	sort.Strings(hosts)

	entries := make([]ReportEntry, 0, len(hosts))
	for _, h := range hosts {
		entries = append(entries, ReportEntry{
			Host:      h,
			OpenPorts: hostPorts[h],
			Changes:   changeCount[h],
			LastSeen:  lastSeen[h],
		})
	}

	return &Report{
		GeneratedAt: time.Now().UTC(),
		Entries:     entries,
	}, nil
}

// SaveReport writes the report as JSON to path.
func SaveReport(path string, r *Report) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create report file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
