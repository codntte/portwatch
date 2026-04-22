package history

import (
	"encoding/json"
	"os"
	"time"
)

// SuppressionRule defines a rule to suppress alerts for a host/port combination.
type SuppressionRule struct {
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// IsExpired returns true if the rule has an expiry and it has passed.
func (r SuppressionRule) IsExpired() bool {
	if r.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(r.ExpiresAt)
}

// SuppressionStore manages suppression rules persisted to disk.
type SuppressionStore struct {
	path string
}

// NewSuppressionStore creates a new SuppressionStore backed by the given file path.
func NewSuppressionStore(path string) *SuppressionStore {
	return &SuppressionStore{path: path}
}

// Add appends a suppression rule, replacing any existing rule for the same host+port.
func (s *SuppressionStore) Add(rule SuppressionRule) error {
	rules, err := s.Load()
	if err != nil {
		return err
	}
	updated := make([]SuppressionRule, 0, len(rules))
	for _, r := range rules {
		if r.Host == rule.Host && r.Port == rule.Port {
			continue
		}
		updated = append(updated, r)
	}
	updated = append(updated, rule)
	return s.save(updated)
}

// Delete removes any suppression rule matching host and port.
func (s *SuppressionStore) Delete(host string, port int) error {
	rules, err := s.Load()
	if err != nil {
		return err
	}
	filtered := make([]SuppressionRule, 0, len(rules))
	for _, r := range rules {
		if r.Host == host && r.Port == port {
			continue
		}
		filtered = append(filtered, r)
	}
	return s.save(filtered)
}

// IsSuppressed returns true if there is a non-expired suppression rule for host+port.
func (s *SuppressionStore) IsSuppressed(host string, port int) (bool, error) {
	rules, err := s.Load()
	if err != nil {
		return false, err
	}
	for _, r := range rules {
		if r.Host == host && r.Port == port && !r.IsExpired() {
			return true, nil
		}
	}
	return false, nil
}

// Load reads all suppression rules from disk.
func (s *SuppressionStore) Load() ([]SuppressionRule, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return []SuppressionRule{}, nil
	}
	if err != nil {
		return nil, err
	}
	var rules []SuppressionRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}
	return rules, nil
}

func (s *SuppressionStore) save(rules []SuppressionRule) error {
	data, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
