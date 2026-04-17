package history

import (
	"encoding/json"
	"os"
	"time"
)

// RetentionPolicy defines how long and how many entries to keep.
type RetentionPolicy struct {
	MaxAge     time.Duration
	MaxEntries int
}

// DefaultRetentionPolicy returns a sensible default policy.
func DefaultRetentionPolicy() RetentionPolicy {
	return RetentionPolicy{
		MaxAge:     7 * 24 * time.Hour,
		MaxEntries: 1000,
	}
}

// Apply loads entries from path, prunes them according to the policy,
// and writes the result back to the same file.
func (p RetentionPolicy) Apply(path string) error {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	var entries []Entry
	for _, line := range splitLines(data) {
		var e Entry
		if json.Unmarshal(line, &e) == nil {
			entries = append(entries, e)
		}
	}

	entries = pruneByAge(entries, p.MaxAge)
	entries = pruneByCount(entries, p.MaxEntries)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for _, e := range entries {
		if err := enc.Encode(e); err != nil {
			return err
		}
	}
	return nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			if i > start {
				lines = append(lines, data[start:i])
			}
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
