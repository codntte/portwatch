package history

import (
	"encoding/json"
	"os"
	"time"
)

// MetricEntry records a numeric measurement for a host at a point in time.
type MetricEntry struct {
	Host      string    `json:"host"`
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// MetricStore persists metric entries to a JSONL file.
type MetricStore struct {
	path string
}

// NewMetricStore returns a MetricStore backed by the given file path.
func NewMetricStore(path string) *MetricStore {
	return &MetricStore{path: path}
}

// Append writes a new MetricEntry to the store.
func (s *MetricStore) Append(e MetricEntry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(e)
}

// Load returns all metric entries from the store.
func (s *MetricStore) Load() ([]MetricEntry, error) {
	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []MetricEntry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e MetricEntry
		if err := dec.Decode(&e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// LoadByHost returns metric entries filtered by host and optional metric name.
func (s *MetricStore) LoadByHost(host, name string) ([]MetricEntry, error) {
	all, err := s.Load()
	if err != nil {
		return nil, err
	}
	var out []MetricEntry
	for _, e := range all {
		if e.Host != host {
			continue
		}
		if name != "" && e.Name != name {
			continue
		}
		out = append(out, e)
	}
	return out, nil
}
