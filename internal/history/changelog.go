// Package history provides persistent storage and querying for portwatch scan history.
// changelog.go tracks a structured changelog of configuration or host-level changes
// made by the user, separate from port scan diffs.
package history

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// ChangeKind describes the category of a changelog entry.
type ChangeKind string

const (
	ChangeKindConfigUpdate ChangeKind = "config_update"
	ChangeKindHostAdded    ChangeKind = "host_added"
	ChangeKindHostRemoved  ChangeKind = "host_removed"
	ChangeKindManual       ChangeKind = "manual"
)

// ChangelogEntry represents a single recorded change event.
type ChangelogEntry struct {
	Timestamp time.Time  `json:"timestamp"`
	Kind      ChangeKind `json:"kind"`
	Host      string     `json:"host,omitempty"`
	Message   string     `json:"message"`
	Author    string     `json:"author,omitempty"`
}

// ChangelogStore manages reading and writing changelog entries to a JSONL file.
type ChangelogStore struct {
	path string
}

// NewChangelogStore returns a ChangelogStore backed by the given file path.
func NewChangelogStore(path string) *ChangelogStore {
	return &ChangelogStore{path: path}
}

// Append adds a new entry to the changelog file.
func (s *ChangelogStore) Append(entry ChangelogEntry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	if entry.Message == "" {
		return errors.New("changelog entry message must not be empty")
	}

	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	line, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = f.Write(append(line, '\n'))
	return err
}

// Load reads all changelog entries from the file.
// Returns an empty slice if the file does not exist.
func (s *ChangelogStore) Load() ([]ChangelogEntry, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []ChangelogEntry{}, nil
		}
		return nil, err
	}

	var entries []ChangelogEntry
	for _, line := range splitLines(string(data)) {
		if line == "" {
			continue
		}
		var e ChangelogEntry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// LoadByKind returns only entries matching the given ChangeKind.
func (s *ChangelogStore) LoadByKind(kind ChangeKind) ([]ChangelogEntry, error) {
	all, err := s.Load()
	if err != nil {
		return nil, err
	}
	var filtered []ChangelogEntry
	for _, e := range all {
		if e.Kind == kind {
			filtered = append(filtered, e)
		}
	}
	return filtered, nil
}
