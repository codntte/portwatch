package runner

import (
	"fmt"
	"net"
	"testing"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/notify"
)

func startTCPServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

// newTestRunner constructs a Runner with a single-host config pointing at the
// given address:port and a stdout dispatcher, reducing boilerplate in tests.
func newTestRunner(t *testing.T, address string, port int) *Runner {
	t.Helper()
	cfg := &config.Config{
		TimeoutSeconds: 1,
		Hosts: []config.Host{
			{Address: address, PortRange: fmt.Sprintf("%d-%d", port, port)},
		},
	}
	ch := notify.NewStdoutChannel()
	d := notify.NewDispatcher([]notify.Channel{ch})
	return New(cfg, d)
}

func TestRun_NoError(t *testing.T) {
	port, stop := startTCPServer(t)
	defer stop()

	r := newTestRunner(t, "127.0.0.1", port)

	if err := r.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_ClosedPort(t *testing.T) {
	// Bind a port and immediately close it so we know it is not listening.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	r := newTestRunner(t, "127.0.0.1", port)

	// Run should succeed even when the target port is closed (closed ports are
	// a normal observation, not an error condition).
	if err := r.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSnapshotFile(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"127.0.0.1", ".portwatch/127_0_0_1.json"},
		{"example.com:8080", ".portwatch/example_com_8080.json"},
	}
	for _, tc := range cases {
		got := snapshotFile(tc.input)
		if got != tc.want {
			t.Errorf("snapshotFile(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
