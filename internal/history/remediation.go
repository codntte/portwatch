package history

import (
	"encoding/json"
	"os"
	"time"
)

// RemediationStatus represents the current state of a remediation action.
type RemediationStatus string

const (
	RemediationPending  RemediationStatus = "pending"
	RemediationApplied  RemediationStatus = "applied"
	RemediationFailed   RemediationStatus = "failed"
	RemediationSkipped  RemediationStatus = "skipped"
)

// RemediationEntry records an automated or manual remediation action for a host/port event.
type RemediationEntry struct {
	ID        string            `json:"id"`
	Host      string            `json:"host"`
	Port      int               `json:"port"`
	Action    string            `json:"action"`
	Status    RemediationStatus `json:"status"`
	Note      string            `json:"note,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// RemediationStore manages persistence of remediation entries.
type RemediationStore struct {
	path string
}

// NewRemediationStore creates a new RemediationStore backed by the given file path.
func NewRemediationStore(path string) *RemediationStore {
	return &RemediationStore{path: path}
}

// Append adds a new remediation entry to the store.
func (s *RemediationStore) Append(e RemediationEntry) error {
	entries, _ := s.Load()
	entries = append(entries, e)
	return s.save(entries)
}

// Load returns all remediation entries from the store.
func (s *RemediationStore) Load() ([]RemediationEntry, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []RemediationEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// UpdateStatus updates the status of a remediation entry by ID.
func (s *RemediationStore) UpdateStatus(id string, status RemediationStatus, note string) error {
	entries, err := s.Load()
	if err != nil {
		return err
	}
	for i, e := range entries {
		if e.ID == id {
			entries[i].Status = status
			entries[i].Note = note
			entries[i].UpdatedAt = time.Now().UTC()
		}
	}
	return s.save(entries)
}

// Delete removes a remediation entry by ID.
func (s *RemediationStore) Delete(id string) error {
	entries, err := s.Load()
	if err != nil {
		return err
	}
	filtered := entries[:0]
	for _, e := range entries {
		if e.ID != id {
			filtered = append(filtered, e)
		}
	}
	return s.save(filtered)
}

func (s *RemediationStore) save(entries []RemediationEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}
