package scanner

import (
	"fmt"
	"strconv"
	"strings"
)

// ParsePortRange parses a port expression into a list of ports.
// Supported formats: "80", "80,443", "8000-8100"
func ParsePortRange(expr string) ([]int, error) {
	var ports []int
	seen := make(map[int]bool)

	parts := strings.Split(expr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			start, err := strconv.Atoi(bounds[0])
			if err != nil {
				return nil, fmt.Errorf("invalid port range start %q", bounds[0])
			}
			end, err := strconv.Atoi(bounds[1])
			if err != nil {
				return nil, fmt.Errorf("invalid port range end %q", bounds[1])
			}
			if start > end || start < 1 || end > 65535 {
				return nil, fmt.Errorf("invalid port range %d-%d", start, end)
			}
			for p := start; p <= end; p++ {
				if !seen[p] {
					ports = append(ports, p)
					seen[p] = true
				}
			}
		} else {
			p, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port %q", part)
			}
			if p < 1 || p > 65535 {
				return nil, fmt.Errorf("port %d out of range", p)
			}
			if !seen[p] {
				ports = append(ports, p)
				seen[p] = true
			}
		}
	}
	return ports, nil
}
