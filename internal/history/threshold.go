package history

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ThresholdRule defines an alert threshold for a host/port combination.
type ThresholdRule struct {
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	MaxClosed int       `json:"max_closed"` // alert if port closed more than N times
	Window    string    `json:"window"`     // duration string e.g. "1h", "24h"
	CreatedAt time.Time `json:"created_at"`
}

// ThresholdStore manages persistence of threshold rules.
type ThresholdStore struct {
	path string
}

// NewThresholdStore returns a ThresholdStore backed by the given file path.
func NewThresholdStore(path string) *ThresholdStore {
	return &ThresholdStore{path: path}
}

// Add appends or replaces a threshold rule for a host+port pair.
func (s *ThresholdStore) Add(rule ThresholdRule) error {
	rules, err := s.Load()
	if err != nil {
		return err
	}
	updated := false
	for i, r := range rules {
		if r.Host == rule.Host && r.Port == rule.Port {
			rules[i] = rule
			updated = true
			break
		}
	}
	if !updated {
		rule.CreatedAt = time.Now().UTC()
		rules = append(rules, rule)
	}
	return s.save(rules)
}

// Delete removes the threshold rule for a host+port pair.
func (s *ThresholdStore) Delete(host string, port int) error {
	rules, err := s.Load()
	if err != nil {
		return err
	}
	filtered := rules[:0]
	for _, r := range rules {
		if r.Host != host || r.Port != port {
			filtered = append(filtered, r)
		}
	}
	return s.save(filtered)
}

// Load reads all threshold rules from disk.
func (s *ThresholdStore) Load() ([]ThresholdRule, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return []ThresholdRule{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("threshold: read: %w", err)
	}
	var rules []ThresholdRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("threshold: unmarshal: %w", err)
	}
	return rules, nil
}

func (s *ThresholdStore) save(rules []ThresholdRule) error {
	data, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return fmt.Errorf("threshold: marshal: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("threshold: write: %w", err)
	}
	return nil
}
