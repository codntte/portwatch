package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"time"
)

// Snapshot holds the open ports observed for a host at a point in time.
type Snapshot struct {
	Host      string    `json:"host"`
	Ports     []int     `json:"ports"`
	ScannedAt time.Time `json:"scanned_at"`
}

// Diff describes ports that changed between two snapshots.
type Diff struct {
	Opened []int
	Closed []int
}

// HasChanges reports whether any ports were opened or closed.
func (d Diff) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

// New creates a Snapshot for the given host and port list.
func New(host string, ports []int) Snapshot {
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)
	return Snapshot{Host: host, Ports: sorted, ScannedAt: time.Now()}
}

// Save writes a snapshot to a JSON file at path.
func Save(path string, s Snapshot) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// Load reads a snapshot from a JSON file at path.
// If the file does not exist, an empty Snapshot is returned with no error.
func Load(path string) (Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Snapshot{}, nil
		}
		return Snapshot{}, err
	}
	defer f.Close()
	var s Snapshot
	return s, json.NewDecoder(f).Decode(&s)
}

// Compare returns ports opened/closed between old and new snapshots.
func Compare(old, new Snapshot) Diff {
	oldSet := toSet(old.Ports)
	newSet := toSet(new.Ports)
	var d Diff
	for p := range newSet {
		if !oldSet[p] {
			d.Opened = append(d.Opened, p)
		}
	}
	for p := range oldSet {
		if !newSet[p] {
			d.Closed = append(d.Closed, p)
		}
	}
	sort.Ints(d.Opened)
	sort.Ints(d.Closed)
	return d
}

func toSet(ports []int) map[int]bool {
	m := make(map[int]bool, len(ports))
	for _, p := range ports {
		m[p] = true
	}
	return m
}
