package runner

import (
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

func TestRun_NoError(t *testing.T) {
	port, stop := startTCPServer(t)
	defer stop()

	cfg := &config.Config{
		TimeoutSeconds: 1,
		Hosts: []config.Host{
			{Address: "127.0.0.1", PortRange: fmt.Sprintf("%d-%d", port, port)},
		},
	}

	ch := notify.NewStdoutChannel()
	d := notify.NewDispatcher([]notify.Channel{ch})
	r := New(cfg, d)

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
