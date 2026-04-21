package history

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// SnapshotDiffEntry records a before/after snapshot comparison for a host.
type SnapshotDiffEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Opened    []int     `json:"opened"`
	Closed    []int     `json:"closed"`
}

// SnapshotDiffStore persists snapshot diff entries to a JSONL file.
type SnapshotDiffStore struct {
	path string
}

// NewSnapshotDiffStore creates a store backed by the given file path.
func NewSnapshotDiffStore(path string) *SnapshotDiffStore {
	return &SnapshotDiffStore{path: path}
}

// Append writes a new SnapshotDiffEntry to the store.
func (s *SnapshotDiffStore) Append(entry SnapshotDiffEntry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("snapshot_diff: open %s: %w", s.path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("snapshot_diff: encode: %w", err)
	}
	return nil
}

// Load reads all SnapshotDiffEntry records from the store.
func (s *SnapshotDiffStore) Load() ([]SnapshotDiffEntry, error) {
	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("snapshot_diff: open %s: %w", s.path, err)
	}
	defer f.Close()
	var entries []SnapshotDiffEntry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e SnapshotDiffEntry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("snapshot_diff: decode: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
