package history

import (
	"encoding/json"
	"os"
	"time"
)

// ProbeResult records the outcome of a single host probe.
type ProbeResult struct {
	Host      string    `json:"host"`
	Timestamp time.Time `json:"timestamp"`
	LatencyMs int64     `json:"latency_ms"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
}

// ProbeStore persists probe results as newline-delimited JSON.
type ProbeStore struct {
	path string
}

// NewProbeStore returns a ProbeStore backed by the given file path.
func NewProbeStore(path string) *ProbeStore {
	return &ProbeStore{path: path}
}

// Append writes a ProbeResult to the store.
func (s *ProbeStore) Append(r ProbeResult) error {
	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	return enc.Encode(r)
}

// Load returns all probe results from the store.
// Returns an empty slice if the file does not exist.
func (s *ProbeStore) Load() ([]ProbeResult, error) {
	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return []ProbeResult{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var results []ProbeResult
	dec := json.NewDecoder(f)
	for dec.More() {
		var r ProbeResult
		if err := dec.Decode(&r); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

// LoadByHost returns probe results filtered to the given host.
func (s *ProbeStore) LoadByHost(host string) ([]ProbeResult, error) {
	all, err := s.Load()
	if err != nil {
		return nil, err
	}
	var out []ProbeResult
	for _, r := range all {
		if r.Host == host {
			out = append(out, r)
		}
	}
	return out, nil
}
