package snapshot

import (
	"encoding/json"
	"os"
	"time"
)

// Snapshot represents the state of scanned ports at a point in time.
type Snapshot struct {
	Host      string         `json:"host"`
	Timestamp time.Time      `json:"timestamp"`
	OpenPorts []int          `json:"open_ports"`
}

// Diff holds the changes between two snapshots.
type Diff struct {
	Opened []int
	Closed []int
}

// New creates a new Snapshot for the given host and open ports.
func New(host string, openPorts []int) *Snapshot {
	return &Snapshot{
		Host:      host,
		Timestamp: time.Now(),
		OpenPorts: openPorts,
	}
}

// Save writes the snapshot to a JSON file at the given path.
func Save(path string, s *Snapshot) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(s)
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

// Compare returns a Diff describing ports opened or closed between prev and curr.
func Compare(prev, curr *Snapshot) Diff {
	prevSet := toSet(prev.OpenPorts)
	currSet := toSet(curr.OpenPorts)

	var d Diff
	for p := range currSet {
		if !prevSet[p] {
			d.Opened = append(d.Opened, p)
		}
	}
	for p := range prevSet {
		if !currSet[p] {
			d.Closed = append(d.Closed, p)
		}
	}
	return d
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
