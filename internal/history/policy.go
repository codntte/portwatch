package history

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// PolicyEntry represents a named scan policy with port range and interval.
type PolicyEntry struct {
	Name      string        `json:"name"`
	Host      string        `json:"host"`
	PortRange string        `json:"port_range"`
	Interval  time.Duration `json:"interval"`
	Enabled   bool          `json:"enabled"`
	CreatedAt time.Time     `json:"created_at"`
}

// PolicyStore manages named scan policies persisted to disk.
type PolicyStore struct {
	path string
}

// NewPolicyStore returns a PolicyStore backed by the given file path.
func NewPolicyStore(path string) *PolicyStore {
	return &PolicyStore{path: path}
}

// Add appends or replaces a policy entry by name.
func (s *PolicyStore) Add(entry PolicyEntry) error {
	entries, err := s.Load()
	if err != nil {
		return err
	}
	updated := make([]PolicyEntry, 0, len(entries))
	for _, e := range entries {
		if e.Name != entry.Name {
			updated = append(updated, e)
		}
	}
	entry.CreatedAt = time.Now()
	updated = append(updated, entry)
	return s.save(updated)
}

// Delete removes a policy by name.
func (s *PolicyStore) Delete(name string) error {
	entries, err := s.Load()
	if err != nil {
		return err
	}
	updated := make([]PolicyEntry, 0, len(entries))
	for _, e := range entries {
		if e.Name != name {
			updated = append(updated, e)
		}
	}
	return s.save(updated)
}

// Get returns a policy by name.
func (s *PolicyStore) Get(name string) (PolicyEntry, bool, error) {
	entries, err := s.Load()
	if err != nil {
		return PolicyEntry{}, false, err
	}
	for _, e := range entries {
		if e.Name == name {
			return e, true, nil
		}
	}
	return PolicyEntry{}, false, nil
}

// Load reads all policy entries from disk.
func (s *PolicyStore) Load() ([]PolicyEntry, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return []PolicyEntry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("policy: read: %w", err)
	}
	var entries []PolicyEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("policy: unmarshal: %w", err)
	}
	return entries, nil
}

func (s *PolicyStore) save(entries []PolicyEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("policy: marshal: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0644); err != nil {
		return fmt.Errorf("policy: write: %w", err)
	}
	return nil
}
