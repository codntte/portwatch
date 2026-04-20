package history

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// Annotation attaches a user-defined note to a specific scan event or diff entry.
type Annotation struct {
	ID        string    `json:"id"`         // arbitrary reference (e.g. timestamp or diff ID)
	Host      string    `json:"host"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

// AnnotationStore manages a JSON file of annotations.
type AnnotationStore struct {
	path string
}

// NewAnnotationStore returns an AnnotationStore backed by the given file path.
func NewAnnotationStore(path string) *AnnotationStore {
	return &AnnotationStore{path: path}
}

// Add appends an annotation to the store.
func (s *AnnotationStore) Add(id, host, note string) error {
	if id == "" || host == "" || note == "" {
		return errors.New("annotation: id, host, and note must not be empty")
	}
	entries, err := s.Load()
	if err != nil {
		return err
	}
	entries = append(entries, Annotation{
		ID:        id,
		Host:      host,
		Note:      note,
		CreatedAt: time.Now().UTC(),
	})
	return s.save(entries)
}

// Delete removes all annotations matching the given id.
func (s *AnnotationStore) Delete(id string) error {
	entries, err := s.Load()
	if err != nil {
		return err
	}
	filtered := entries[:0]
	for _, a := range entries {
		if a.ID != id {
			filtered = append(filtered, a)
		}
	}
	return s.save(filtered)
}

// Load reads all annotations from disk. Returns an empty slice if the file
// does not exist.
func (s *AnnotationStore) Load() ([]Annotation, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return []Annotation{}, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []Annotation
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (s *AnnotationStore) save(entries []Annotation) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
