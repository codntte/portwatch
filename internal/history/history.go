package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single scan event recorded in history.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Opened    []int     `json:"opened"`
	Closed    []int     `json:"closed"`
}

// Store manages persisted scan history for a host.
type Store struct {
	dir string
}

// NewStore creates a Store that writes history files under dir.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

// Append adds an entry to the host's history file.
func (s *Store) Append(e Entry) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return fmt.Errorf("history: mkdir: %w", err)
	}
	f, err := os.OpenFile(s.filePath(e.Host), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("history: open: %w", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(e)
}

// Load returns all entries for the given host.
func (s *Store) Load(host string) ([]Entry, error) {
	f, err := os.Open(s.filePath(host))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return fmt.Errorf("history: open: %w", err)
	}
	defer f.Close()
	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("history: decode: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (s *Store) filePath(host string) string {
	safe := filepath.Base(host)
	return filepath.Join(s.dir, safe+".jsonl")
}
