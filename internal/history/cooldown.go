package history

import (
	"encoding/json"
	"os"
	"time"
)

// CooldownEntry represents a suppression cooldown for a host+port combination.
type CooldownEntry struct {
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Until     time.Time `json:"until"`
	CreatedAt time.Time `json:"created_at"`
}

// CooldownStore manages alert cooldowns to prevent alert fatigue.
type CooldownStore struct {
	path string
}

// NewCooldownStore returns a CooldownStore backed by the given file path.
func NewCooldownStore(path string) *CooldownStore {
	return &CooldownStore{path: path}
}

// Add writes a cooldown entry for the given host and port, active until the given time.
func (s *CooldownStore) Add(host string, port int, until time.Time) error {
	entries, _ := s.Load()
	// Replace existing entry for same host+port.
	updated := entries[:0]
	for _, e := range entries {
		if e.Host != host || e.Port != port {
			updated = append(updated, e)
		}
	}
	updated = append(updated, CooldownEntry{
		Host:      host,
		Port:      port,
		Until:     until,
		CreatedAt: time.Now().UTC(),
	})
	return s.save(updated)
}

// IsActive returns true if there is an active cooldown for the given host and port.
func (s *CooldownStore) IsActive(host string, port int) bool {
	entries, err := s.Load()
	if err != nil {
		return false
	}
	now := time.Now().UTC()
	for _, e := range entries {
		if e.Host == host && e.Port == port && now.Before(e.Until) {
			return true
		}
	}
	return false
}

// Delete removes the cooldown entry for the given host and port.
func (s *CooldownStore) Delete(host string, port int) error {
	entries, _ := s.Load()
	filtered := entries[:0]
	for _, e := range entries {
		if e.Host != host || e.Port != port {
			filtered = append(filtered, e)
		}
	}
	return s.save(filtered)
}

// Load reads all cooldown entries from disk.
func (s *CooldownStore) Load() ([]CooldownEntry, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []CooldownEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (s *CooldownStore) save(entries []CooldownEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
