package history

import (
	"encoding/json"
	"os"
	"time"
)

// IncidentSeverity represents the severity level of a port incident.
type IncidentSeverity string

const (
	SeverityLow    IncidentSeverity = "low"
	SeverityMedium IncidentSeverity = "medium"
	SeverityHigh   IncidentSeverity = "high"
)

// Incident represents a detected port change event that has been escalated.
type Incident struct {
	ID        string           `json:"id"`
	Host      string           `json:"host"`
	Ports     []int            `json:"ports"`
	Kind      string           `json:"kind"` // "opened" or "closed"
	Severity  IncidentSeverity `json:"severity"`
	Note      string           `json:"note,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	ResolvedAt *time.Time      `json:"resolved_at,omitempty"`
}

// IncidentStore manages persistence of incidents.
type IncidentStore struct {
	path string
}

// NewIncidentStore returns a new IncidentStore backed by the given file path.
func NewIncidentStore(path string) *IncidentStore {
	return &IncidentStore{path: path}
}

// Append adds a new incident to the store.
func (s *IncidentStore) Append(inc Incident) error {
	incidents, err := s.Load()
	if err != nil {
		return err
	}
	incidents = append(incidents, inc)
	return s.save(incidents)
}

// Load reads all incidents from disk.
func (s *IncidentStore) Load() ([]Incident, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return []Incident{}, nil
	}
	if err != nil {
		return nil, err
	}
	var incidents []Incident
	if err := json.Unmarshal(data, &incidents); err != nil {
		return nil, err
	}
	return incidents, nil
}

// Resolve marks an incident as resolved by ID.
func (s *IncidentStore) Resolve(id string) error {
	incidents, err := s.Load()
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	for i, inc := range incidents {
		if inc.ID == id && inc.ResolvedAt == nil {
			incidents[i].ResolvedAt = &now
		}
	}
	return s.save(incidents)
}

func (s *IncidentStore) save(incidents []Incident) error {
	data, err := json.MarshalIndent(incidents, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}
