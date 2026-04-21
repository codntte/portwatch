package history

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// WatchlistEntry represents a host+port combination being actively watched.
type WatchlistEntry struct {
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Label     string    `json:"label,omitempty"`
	AddedAt   time.Time `json:"added_at"`
}

// WatchlistStore manages a persisted list of watched host:port pairs.
type WatchlistStore struct {
	path string
}

// NewWatchlistStore creates a new WatchlistStore backed by the given file path.
func NewWatchlistStore(path string) *WatchlistStore {
	return &WatchlistStore{path: path}
}

// Add appends a new entry to the watchlist. Duplicate host+port pairs are ignored.
func (s *WatchlistStore) Add(host string, port int, label string) error {
	entries, err := s.Load()
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.Host == host && e.Port == port {
			return nil
		}
	}
	entries = append(entries, WatchlistEntry{
		Host:    host,
		Port:    port,
		Label:   label,
		AddedAt: time.Now().UTC(),
	})
	return s.save(entries)
}

// Delete removes the entry matching host+port from the watchlist.
func (s *WatchlistStore) Delete(host string, port int) error {
	entries, err := s.Load()
	if err != nil {
		return err
	}
	filtered := entries[:0]
	for _, e := range entries {
		if e.Host != host || e.Port != port {
			filtered = append(filtered, e)
		}
	}
	return s.save(filtered)
}

// Load reads all watchlist entries from disk.
func (s *WatchlistStore) Load() ([]WatchlistEntry, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return []WatchlistEntry{}, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []WatchlistEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (s *WatchlistStore) save(entries []WatchlistEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
