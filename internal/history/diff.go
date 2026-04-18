package history

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// DiffEntry records a port change event between two snapshots.
type DiffEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Opened    []int     `json:"opened"`
	Closed    []int     `json:"closed"`
}

// AppendDiff appends a DiffEntry to the given file (one JSON object per line).
func AppendDiff(path string, entry DiffEntry) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("history/diff: open %s: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("history/diff: encode: %w", err)
	}
	return nil
}

// LoadDiffs reads all DiffEntry records from the given file.
func LoadDiffs(path string) ([]DiffEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("history/diff: open %s: %w", path, err)
	}
	defer f.Close()

	var entries []DiffEntry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e DiffEntry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("history/diff: decode: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
