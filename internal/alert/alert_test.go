package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/snapshot"
)

func TestNotify_OpenedAndClosed(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New()
	n.Out = &buf

	diff := snapshot.Diff{
		Opened: []int{8080, 443},
		Closed: []int{22},
	}

	n.Notify("localhost", diff)

	out := buf.String()

	if !strings.Contains(out, "port 8080 newly OPEN") {
		t.Errorf("expected open alert for 8080, got:\n%s", out)
	}
	if !strings.Contains(out, "port 443 newly OPEN") {
		t.Errorf("expected open alert for 443, got:\n%s", out)
	}
	if !strings.Contains(out, "port 22 newly CLOSED") {
		t.Errorf("expected closed alert for 22, got:\n%s", out)
	}
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT level in output")
	}
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN level in output")
	}
}

func TestNotify_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New()
	n.Out = &buf

	n.Notify("localhost", snapshot.Diff{})

	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got: %s", buf.String())
	}
}
