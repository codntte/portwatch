package history

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// Label associates a named label with a host for grouping/filtering.
type Label struct {
	Host      string    `json:"host"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// LabelStore manages persistent labels for hosts.
type LabelStore struct {
	path string
}

// NewLabelStore returns a LabelStore backed by the given file path.
func NewLabelStore(path string) *LabelStore {
	return &LabelStore{path: path}
}

// Add appends a label for the given host. Duplicate host+name pairs are ignored.
func (s *LabelStore) Add(host, name string) error {
	if host == "" || name == "" {
		return errors.New("host and name must not be empty")
	}
	labels, err := s.Load("")
	if err != nil {
		return err
	}
	for _, l := range labels {
		if l.Host == host && l.Name == name {
			return nil
		}
	}
	labels = append(labels, Label{Host: host, Name: name, CreatedAt: time.Now().UTC()})
	return s.save(labels)
}

// Delete removes all labels matching the given host and name.
func (s *LabelStore) Delete(host, name string) error {
	labels, err := s.Load("")
	if err != nil {
		return err
	}
	filtered := labels[:0]
	for _, l := range labels {
		if !(l.Host == host && l.Name == name) {
			filtered = append(filtered, l)
		}
	}
	return s.save(filtered)
}

// Load returns all labels, optionally filtered by host (empty string = all).
func (s *LabelStore) Load(host string) ([]Label, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return []Label{}, nil
	}
	if err != nil {
		return nil, err
	}
	var labels []Label
	if err := json.Unmarshal(data, &labels); err != nil {
		return nil, err
	}
	if host == "" {
		return labels, nil
	}
	var out []Label
	for _, l := range labels {
		if l.Host == host {
			out = append(out, l)
		}
	}
	return out, nil
}

func (s *LabelStore) save(labels []Label) error {
	data, err := json.MarshalIndent(labels, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
