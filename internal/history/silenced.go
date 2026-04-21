package history

import (
	"encoding/json"
	"os"
	"time"
)

// SilencedEntry represents a suppressed alert rule for a host/port combination.
type SilencedEntry struct {
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// IsExpired reports whether the silence period has passed.
func (e SilencedEntry) IsExpired() bool {
	if e.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(e.ExpiresAt)
}

// SilencedStore manages persisted silence rules.
type SilencedStore struct {
	path string
}

// NewSilencedStore creates a SilencedStore backed by the given file path.
func NewSilencedStore(path string) *SilencedStore {
	return &SilencedStore{path: path}
}

// Add appends a new silence entry, replacing any existing entry for the same host+port.
func (s *SilencedStore) Add(entry SilencedEntry) error {
	entries, _ := s.Load()
	updated := make([]SilencedEntry, 0, len(entries)+1)
	for _, e := range entries {
		if e.Host == entry.Host && e.Port == entry.Port {
			continue
		}
		updated = append(updated, e)
	}
	updated = append(updated, entry)
	return s.save(updated)
}

// Delete removes the silence rule for the given host and port.
func (s *SilencedStore) Delete(host string, port int) error {
	entries, _ := s.Load()
	updated := make([]SilencedEntry, 0, len(entries))
	for _, e := range entries {
		if e.Host == host && e.Port == port {
			continue
		}
		updated = append(updated, e)
	}
	return s.save(updated)
}

// Load returns all non-expired silence entries.
func (s *SilencedStore) Load() ([]SilencedEntry, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []SilencedEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	active := entries[:0]
	for _, e := range entries {
		if !e.IsExpired() {
			active = append(active, e)
		}
	}
	return active, nil
}

// IsSilenced reports whether the given host+port has an active silence rule.
func (s *SilencedStore) IsSilenced(host string, port int) bool {
	entries, _ := s.Load()
	for _, e := range entries {
		if e.Host == host && e.Port == port {
			return true
		}
	}
	return false
}

func (s *SilencedStore) save(entries []SilencedEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
