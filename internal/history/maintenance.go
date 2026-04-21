package history

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// MaintenanceWindow represents a scheduled maintenance period for a host.
type MaintenanceWindow struct {
	ID        string    `json:"id"`
	Host      string    `json:"host"`
	StartsAt  time.Time `json:"starts_at"`
	EndsAt    time.Time `json:"ends_at"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

// IsActive reports whether the window is currently active.
func (m MaintenanceWindow) IsActive() bool {
	now := time.Now()
	return now.After(m.StartsAt) && now.Before(m.EndsAt)
}

// NewMaintenanceStore creates a MaintenanceStore backed by path.
func NewMaintenanceStore(path string) *MaintenanceStore {
	return &MaintenanceStore{path: path}
}

// MaintenanceStore persists maintenance windows as newline-delimited JSON.
type MaintenanceStore struct {
	path string
}

// Add appends a new maintenance window to the store.
func (s *MaintenanceStore) Add(w MaintenanceWindow) error {
	f, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("maintenance: open %s: %w", s.path, err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(w)
}

// Load returns all stored maintenance windows.
func (s *MaintenanceStore) Load() ([]MaintenanceWindow, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("maintenance: read %s: %w", s.path, err)
	}
	return parseMaintenanceLines(data)
}

// Delete removes the window with the given ID.
func (s *MaintenanceStore) Delete(id string) error {
	windows, err := s.Load()
	if err != nil {
		return err
	}
	filtered := windows[:0]
	for _, w := range windows {
		if w.ID != id {
			filtered = append(filtered, w)
		}
	}
	return rewriteMaintenanceFile(s.path, filtered)
}

// ActiveFor returns windows currently active for the given host.
func (s *MaintenanceStore) ActiveFor(host string) ([]MaintenanceWindow, error) {
	all, err := s.Load()
	if err != nil {
		return nil, err
	}
	var active []MaintenanceWindow
	for _, w := range all {
		if w.Host == host && w.IsActive() {
			active = append(active, w)
		}
	}
	return active, nil
}

func parseMaintenanceLines(data []byte) ([]MaintenanceWindow, error) {
	var windows []MaintenanceWindow
	for _, line := range splitLines(data) {
		var w MaintenanceWindow
		if err := json.Unmarshal([]byte(line), &w); err != nil {
			return nil, fmt.Errorf("maintenance: parse: %w", err)
		}
		windows = append(windows, w)
	}
	return windows, nil
}

func rewriteMaintenanceFile(path string, windows []MaintenanceWindow) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("maintenance: rewrite %s: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, w := range windows {
		if err := enc.Encode(w); err != nil {
			return err
		}
	}
	return nil
}
