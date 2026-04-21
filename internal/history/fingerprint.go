package history

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// FingerprintEntry records a SHA-256 fingerprint of a host's open port set at a point in time.
type FingerprintEntry struct {
	Host      string    `json:"host"`
	Ports     []int     `json:"ports"`
	Fingerprint string  `json:"fingerprint"`
	Timestamp time.Time `json:"timestamp"`
}

// ComputeFingerprint returns a deterministic SHA-256 hex digest for a sorted port list.
func ComputeFingerprint(host string, ports []int) string {
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)

	h := sha256.New()
	fmt.Fprintf(h, "%s:", host)
	for _, p := range sorted {
		fmt.Fprintf(h, "%d,", p)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// NewFingerprintStore creates a FingerprintStore backed by the given file path.
func NewFingerprintStore(path string) *FingerprintStore {
	return &FingerprintStore{path: path}
}

// FingerprintStore persists fingerprint entries as newline-delimited JSON.
type FingerprintStore struct {
	path string
}

// Append adds a new fingerprint entry to the store.
func (s *FingerprintStore) Append(host string, ports []int) error {
	entry := FingerprintEntry{
		Host:        host,
		Ports:       ports,
		Fingerprint: ComputeFingerprint(host, ports),
		Timestamp:   time.Now().UTC(),
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("fingerprint marshal: %w", err)
	}
	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("fingerprint open: %w", err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s\n", data)
	return err
}

// Load returns all fingerprint entries from the store.
func (s *FingerprintStore) Load() ([]FingerprintEntry, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fingerprint read: %w", err)
	}
	var entries []FingerprintEntry
	for _, line := range splitLines(string(data)) {
		var e FingerprintEntry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			continue
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// LatestByHost returns the most recent fingerprint entry per host.
func (s *FingerprintStore) LatestByHost() (map[string]FingerprintEntry, error) {
	entries, err := s.Load()
	if err != nil {
		return nil, err
	}
	result := make(map[string]FingerprintEntry)
	for _, e := range entries {
		if prev, ok := result[e.Host]; !ok || e.Timestamp.After(prev.Timestamp) {
			result[e.Host] = e
		}
	}
	return result, nil
}
