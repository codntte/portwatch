package history

import (
	"encoding/json"
	"os"
	"time"
)

// Acknowledgment represents a user acknowledgment of a port change or alert.
type Acknowledgment struct {
	ID        string    `json:"id"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	AckedBy   string    `json:"acked_by"`
	Comment   string    `json:"comment"`
	AckedAt   time.Time `json:"acked_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// AcknowledgmentStore manages persistence of acknowledgments.
type AcknowledgmentStore struct {
	path string
}

// NewAcknowledgmentStore creates a new AcknowledgmentStore backed by the given file.
func NewAcknowledgmentStore(path string) *AcknowledgmentStore {
	return &AcknowledgmentStore{path: path}
}

// Append adds an acknowledgment entry to the store.
func (s *AcknowledgmentStore) Append(a Acknowledgment) error {
	entries, _ := s.Load()
	entries = append(entries, a)
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Load returns all acknowledgment entries from the store.
func (s *AcknowledgmentStore) Load() ([]Acknowledgment, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []Acknowledgment
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// LoadByHost returns acknowledgments filtered by host.
func (s *AcknowledgmentStore) LoadByHost(host string) ([]Acknowledgment, error) {
	all, err := s.Load()
	if err != nil {
		return nil, err
	}
	var result []Acknowledgment
	for _, a := range all {
		if a.Host == host {
			result = append(result, a)
		}
	}
	return result, nil
}

// Delete removes all acknowledgments matching the given host and port.
func (s *AcknowledgmentStore) Delete(host string, port int) error {
	all, err := s.Load()
	if err != nil {
		return err
	}
	var kept []Acknowledgment
	for _, a := range all {
		if a.Host != host || a.Port != port {
			kept = append(kept, a)
		}
	}
	data, err := json.MarshalIndent(kept, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
