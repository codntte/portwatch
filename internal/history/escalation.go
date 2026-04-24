package history

import (
	"encoding/json"
	"os"
	"time"
)

// EscalationLevel represents the severity of an escalation.
type EscalationLevel string

const (
	EscalationInfo     EscalationLevel = "info"
	EscalationWarning  EscalationLevel = "warning"
	EscalationCritical EscalationLevel = "critical"
)

// EscalationEntry records an escalation event for a host.
type EscalationEntry struct {
	Host      string          `json:"host"`
	Level     EscalationLevel `json:"level"`
	Reason    string          `json:"reason"`
	Triggered time.Time       `json:"triggered"`
	Resolved  *time.Time      `json:"resolved,omitempty"`
}

// EscalationStore manages persistence of escalation entries.
type EscalationStore struct {
	path string
}

// NewEscalationStore creates a new EscalationStore backed by the given file path.
func NewEscalationStore(path string) *EscalationStore {
	return &EscalationStore{path: path}
}

// Append adds a new escalation entry to the store.
func (s *EscalationStore) Append(e EscalationEntry) error {
	entries, err := s.Load()
	if err != nil {
		return err
	}
	entries = append(entries, e)
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

// Load reads all escalation entries from disk.
func (s *EscalationStore) Load() ([]EscalationEntry, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return []EscalationEntry{}, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []EscalationEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// Resolve marks the first unresolved entry for the given host as resolved.
func (s *EscalationStore) Resolve(host string, at time.Time) error {
	entries, err := s.Load()
	if err != nil {
		return err
	}
	for i := range entries {
		if entries[i].Host == host && entries[i].Resolved == nil {
			entries[i].Resolved = &at
			break
		}
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

// LoadByHost returns all escalation entries for a specific host.
func (s *EscalationStore) LoadByHost(host string) ([]EscalationEntry, error) {
	all, err := s.Load()
	if err != nil {
		return nil, err
	}
	var result []EscalationEntry
	for _, e := range all {
		if e.Host == host {
			result = append(result, e)
		}
	}
	return result, nil
}
