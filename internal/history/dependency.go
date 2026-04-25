package history

import (
	"encoding/json"
	"os"
	"time"
)

// Dependency represents a tracked relationship between two hosts,
// indicating that one host depends on another for connectivity.
type Dependency struct {
	ID        string    `json:"id"`
	Host      string    `json:"host"`
	DependsOn string    `json:"depends_on"`
	Ports     []int     `json:"ports,omitempty"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// DependencyStore manages host dependency records.
type DependencyStore struct {
	path string
}

// NewDependencyStore returns a DependencyStore backed by the given file path.
func NewDependencyStore(path string) *DependencyStore {
	return &DependencyStore{path: path}
}

// Add appends a new dependency entry to the store.
func (s *DependencyStore) Add(dep Dependency) error {
	deps, err := s.Load()
	if err != nil {
		return err
	}
	dep.CreatedAt = time.Now().UTC()
	deps = append(deps, dep)
	return s.save(deps)
}

// Load reads all dependency entries from the store.
func (s *DependencyStore) Load() ([]Dependency, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return []Dependency{}, nil
	}
	if err != nil {
		return nil, err
	}
	var deps []Dependency
	if err := json.Unmarshal(data, &deps); err != nil {
		return nil, err
	}
	return deps, nil
}

// Delete removes a dependency by ID.
func (s *DependencyStore) Delete(id string) error {
	deps, err := s.Load()
	if err != nil {
		return err
	}
	filtered := deps[:0]
	for _, d := range deps {
		if d.ID != id {
			filtered = append(filtered, d)
		}
	}
	return s.save(filtered)
}

// LoadByHost returns all dependencies where Host or DependsOn matches the given host.
func (s *DependencyStore) LoadByHost(host string) ([]Dependency, error) {
	all, err := s.Load()
	if err != nil {
		return nil, err
	}
	var result []Dependency
	for _, d := range all {
		if d.Host == host || d.DependsOn == host {
			result = append(result, d)
		}
	}
	return result, nil
}

func (s *DependencyStore) save(deps []Dependency) error {
	data, err := json.MarshalIndent(deps, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}
