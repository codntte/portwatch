package history

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// HostStatus represents the last known status of a host.
type HostStatus struct {
	Host      string    `json:"host"`
	OpenPorts []int     `json:"open_ports"`
	LastSeen  time.Time `json:"last_seen"`
	Up        bool      `json:"up"`
}

// HostStatusStore manages persisted host status records.
type HostStatusStore struct {
	path string
}

// NewHostStatusStore creates a new HostStatusStore backed by the given file.
func NewHostStatusStore(path string) *HostStatusStore {
	return &HostStatusStore{path: path}
}

// Upsert inserts or replaces the status for a host.
func (s *HostStatusStore) Upsert(status HostStatus) error {
	records, err := s.LoadAll()
	if err != nil {
		return err
	}
	updated := false
	for i, r := range records {
		if r.Host == status.Host {
			records[i] = status
			updated = true
			break
		}
	}
	if !updated {
		records = append(records, status)
	}
	return s.save(records)
}

// Get returns the status for a specific host, or an error if not found.
func (s *HostStatusStore) Get(host string) (HostStatus, error) {
	records, err := s.LoadAll()
	if err != nil {
		return HostStatus{}, err
	}
	for _, r := range records {
		if r.Host == host {
			return r, nil
		}
	}
	return HostStatus{}, fmt.Errorf("host %q not found", host)
}

// LoadAll returns all host status records sorted by host name.
func (s *HostStatusStore) LoadAll() ([]HostStatus, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return []HostStatus{}, nil
	}
	if err != nil {
		return nil, err
	}
	var records []HostStatus
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Host < records[j].Host
	})
	return records, nil
}

func (s *HostStatusStore) save(records []HostStatus) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
