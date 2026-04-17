package history

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrintTags_Basic(t *testing.T) {
	dir := t.TempDir()
	s := NewTagStore(filepath.Join(dir, "tags.json"))
	_ = s.Add("release-1", "first")
	_ = s.Add("release-2", "")

	var buf bytes.Buffer
	if err := PrintTags(s, &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "release-1") {
		t.Error("expected release-1 in output")
	}
	if !strings.Contains(out, "first") {
		t.Error("expected note 'first' in output")
	}
	if !strings.Contains(out, "release-2") {
		t.Error("expected release-2 in output")
	}
}

func TestPrintTags_NoTags(t *testing.T) {
	dir := t.TempDir()
	s := NewTagStore(filepath.Join(dir, "tags.json"))

	var buf bytes.Buffer
	if err := PrintTags(s, &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "no tags") {
		t.Error("expected 'no tags' message")
	}
}
