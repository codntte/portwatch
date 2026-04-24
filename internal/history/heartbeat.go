package history

import (
	"encoding/json"
	"os"
	"time"
)

// HeartbeatEntry records a periodic liveness ping for a host.
type HeartbeatEntry struct {
	Host      string    `json:"host"`
	Timestamp time.Time `json:"timestamp"`
	LatencyMs int64     `json:"latency_ms"`
	Alive     bool      `json:"alive"`
}

// HeartbeatStore manages persisted heartbeat records.
type HeartbeatStore struct {
	path string
}

// NewHeartbeatStore returns a HeartbeatStore backed by the given file path.
func NewHeartbeatStore(path string) *HeartbeatStore {
	return &HeartbeatStore{path: path}
}

// Append adds a new heartbeat entry to the store.
func (s *HeartbeatStore) Append(entry HeartbeatEntry) error {
	f, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(entry)
}

// Load returns all heartbeat entries from the store.
func (s *HeartbeatStore) Load() ([]HeartbeatEntry, error) {
	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []HeartbeatEntry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e HeartbeatEntry
		if err := dec.Decode(&e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// LoadByHost returns heartbeat entries filtered to a specific host.
func (s *HeartbeatStore) LoadByHost(host string) ([]HeartbeatEntry, error) {
	all, err := s.Load()
	if err != nil {
		return nil, err
	}
	var out []HeartbeatEntry
	for _, e := range all {
		if e.Host == host {
			out = append(out, e)
		}
	}
	return out, nil
}

// LastSeen returns the most recent heartbeat entry for a host, or nil if none.
func (s *HeartbeatStore) LastSeen(host string) (*HeartbeatEntry, error) {
	entries, err := s.LoadByHost(host)
	if err != nil || len(entries) == 0 {
		return nil, err
	}
	last := entries[len(entries)-1]
	return &last, nil
}
