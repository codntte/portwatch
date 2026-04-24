package history

import (
	"encoding/json"
	"os"
	"time"
)

// AuditAction represents the type of action recorded in the audit log.
type AuditAction string

const (
	AuditActionAdd    AuditAction = "add"
	AuditActionDelete AuditAction = "delete"
	AuditActionUpdate AuditAction = "update"
)

// AuditEntry represents a single audit log record.
type AuditEntry struct {
	Timestamp time.Time   `json:"timestamp"`
	Actor     string      `json:"actor"`
	Action    AuditAction `json:"action"`
	Target    string      `json:"target"`
	Detail    string      `json:"detail,omitempty"`
}

// AuditStore manages appending and loading audit log entries.
type AuditStore struct {
	path string
}

// NewAuditStore creates an AuditStore backed by the given file path.
func NewAuditStore(path string) *AuditStore {
	return &AuditStore{path: path}
}

// Append writes a new audit entry to the log file.
func (s *AuditStore) Append(entry AuditEntry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(entry)
}

// Load returns all audit entries from the log file.
func (s *AuditStore) Load() ([]AuditEntry, error) {
	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []AuditEntry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e AuditEntry
		if err := dec.Decode(&e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// LoadByActor returns entries matching the given actor name.
func (s *AuditStore) LoadByActor(actor string) ([]AuditEntry, error) {
	all, err := s.Load()
	if err != nil {
		return nil, err
	}
	var out []AuditEntry
	for _, e := range all {
		if e.Actor == actor {
			out = append(out, e)
		}
	}
	return out, nil
}
