package notify

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestStdoutChannel_Send(t *testing.T) {
	var buf bytes.Buffer
	ch := &StdoutChannel{Writer: &buf}

	err := ch.Send("test subject", "test body")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "test subject") {
		t.Errorf("expected subject in output, got: %s", out)
	}
	if !strings.Contains(out, "test body") {
		t.Errorf("expected body in output, got: %s", out)
	}
}

func TestDispatcher_Dispatch_WithChanges(t *testing.T) {
	var buf bytes.Buffer
	ch := &StdoutChannel{Writer: &buf}
	d := NewDispatcher(ch)

	diff := alert.Diff{
		Opened: []int{80, 443},
		Closed: []int{8080},
	}

	err := d.Dispatch("localhost", diff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected host in output")
	}
	if !strings.Contains(out, "OPENED") {
		t.Errorf("expected OPENED in output")
	}
	if !strings.Contains(out, "CLOSED") {
		t.Errorf("expected CLOSED in output")
	}
}

func TestDispatcher_Dispatch_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	ch := &StdoutChannel{Writer: &buf}
	d := NewDispatcher(ch)

	diff := alert.Diff{}
	err := d.Dispatch("localhost", diff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got: %s", buf.String())
	}
}
