package history

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// Baseline represents a named snapshot of open ports for a host at a point in time.
type Baseline struct {
	Name      string         `json:"name"`
	Host      string         `json:"host"`
	Ports     []int          `json:"ports"`
	CreatedAt time.Time      `json:"created_at"`
}

// BaselineStore manages named baselines persisted to a JSON file.
type BaselineStore struct {
	path string
}

// NewBaselineStore creates a BaselineStore backed by the given file path.
func NewBaselineStore(path string) *BaselineStore {
	return &BaselineStore{path: path}
}

func (s *BaselineStore) load() ([]Baseline, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return []Baseline{}, nil
	}
	if err != nil {
		return nil, err
	}
	var baselines []Baseline
	if err := json.Unmarshal(data, &baselines); err != nil {
		return nil, err
	}
	return baselines, nil
}

func (s *BaselineStore) save(baselines []Baseline) error {
	data, err := json.MarshalIndent(baselines, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Add saves or replaces a baseline with the given name and host.
func (s *BaselineStore) Add(name, host string, ports []int) error {
	baselines, err := s.load()
	if err != nil {
		return err
	}
	for i, b := range baselines {
		if b.Name == name && b.Host == host {
			baselines[i] = Baseline{Name: name, Host: host, Ports: ports, CreatedAt: time.Now().UTC()}
			return s.save(baselines)
		}
	}
	baselines = append(baselines, Baseline{
		Name:      name,
		Host:      host,
		Ports:     ports,
		CreatedAt: time.Now().UTC(),
	})
	return s.save(baselines)
}

// Get returns the baseline with the given name and host, or false if not found.
func (s *BaselineStore) Get(name, host string) (Baseline, bool, error) {
	baselines, err := s.load()
	if err != nil {
		return Baseline{}, false, err
	}
	for _, b := range baselines {
		if b.Name == name && b.Host == host {
			return b, true, nil
		}
	}
	return Baseline{}, false, nil
}

// Delete removes the baseline with the given name and host.
func (s *BaselineStore) Delete(name, host string) error {
	baselines, err := s.load()
	if err != nil {
		return err
	}
	filtered := baselines[:0]
	for _, b := range baselines {
		if b.Name != name || b.Host != host {
			filtered = append(filtered, b)
		}
	}
	return s.save(filtered)
}

// List returns all stored baselines.
func (s *BaselineStore) List() ([]Baseline, error) {
	return s.load()
}
