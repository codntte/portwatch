package scanner

import (
	"net"
	"strconv"
	"testing"
	"time"
)

// startTCPServer starts a local TCP listener and returns its port.
func startTCPServer(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	return port
}

func TestScanPort_Open(t *testing.T) {
	port := startTCPServer(t)
	state := ScanPort("127.0.0.1", port, 2*time.Second)
	if !state.Open {
		t.Errorf("expected port %d to be open", port)
	}
}

func TestScanPort_Closed(t *testing.T) {
	state := ScanPort("127.0.0.1", 1, 500*time.Millisecond)
	if state.Open {
		t.Errorf("expected port 1 to be closed")
	}
}

func TestScanPorts(t *testing.T) {
	port := startTCPServer(t)
	opts := DefaultOptions()
	results := ScanPorts("127.0.0.1", []int{port, 1}, opts)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	openCount := 0
	for _, r := range results {
		if r.Open {
			openCount++
		}
	}
	if openCount != 1 {
		t.Errorf("expected 1 open port, got %d", openCount)
	}
}
