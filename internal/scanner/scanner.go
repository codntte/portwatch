package scanner

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a scanned port.
type PortState struct {
	Host   string
	Port   int
	Open   bool
	Latency time.Duration
}

// ScanOptions holds configuration for a port scan.
type ScanOptions struct {
	Timeout    time.Duration
	Concurrency int
}

// DefaultOptions returns sensible defaults for scanning.
func DefaultOptions() ScanOptions {
	return ScanOptions{
		Timeout:     2 * time.Second,
		Concurrency: 100,
	}
}

// ScanPort checks whether a single TCP port is open on the given host.
func ScanPort(host string, port int, timeout time.Duration) PortState {
	address := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, timeout)
	latency := time.Since(start)
	if err != nil {
		return PortState{Host: host, Port: port, Open: false}
	}
	conn.Close()
	return PortState{Host: host, Port: port, Open: true, Latency: latency}
}

// ScanPorts scans a range of ports concurrently and returns results.
func ScanPorts(host string, ports []int, opts ScanOptions) []PortState {
	results := make([]PortState, 0, len(ports))
	sem := make(chan struct{}, opts.Concurrency)
	resultCh := make(chan PortState, len(ports))

	for _, port := range ports {
		sem <- struct{}{}
		go func(p int) {
			defer func() { <-sem }()
			resultCh <- ScanPort(host, p, opts.Timeout)
		}(port)
	}

	for range ports {
		results = append(results, <-resultCh)
	}
	return results
}
