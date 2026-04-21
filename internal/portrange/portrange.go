// Package portrange provides utilities for parsing and expanding
// port range expressions such as "80", "8080-8090", or "22,80,443".
package portrange

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	minPort = 1
	maxPort = 65535
)

// Parse accepts a comma-separated list of port numbers and port ranges
// (e.g. "22,80,8080-8090") and returns a deduplicated, sorted slice of
// individual port numbers. An error is returned if any token is invalid.
func Parse(expr string) ([]int, error) {
	if strings.TrimSpace(expr) == "" {
		return nil, fmt.Errorf("portrange: empty expression")
	}

	seen := make(map[int]struct{})
	tokens := strings.Split(expr, ",")

	for _, tok := range tokens {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}

		if strings.Contains(tok, "-") {
			parts := strings.SplitN(tok, "-", 2)
			lo, err := parsePort(parts[0])
			if err != nil {
				return nil, fmt.Errorf("portrange: invalid range start %q: %w", parts[0], err)
			}
			hi, err := parsePort(parts[1])
			if err != nil {
				return nil, fmt.Errorf("portrange: invalid range end %q: %w", parts[1], err)
			}
			if lo > hi {
				return nil, fmt.Errorf("portrange: range start %d is greater than end %d", lo, hi)
			}
			for p := lo; p <= hi; p++ {
				seen[p] = struct{}{}
			}
		} else {
			p, err := parsePort(tok)
			if err != nil {
				return nil, fmt.Errorf("portrange: invalid port %q: %w", tok, err)
			}
			seen[p] = struct{}{}
		}
	}

	return sorted(seen), nil
}

// Contains reports whether port is within any range described by expr.
func Contains(expr string, port int) (bool, error) {
	ports, err := Parse(expr)
	if err != nil {
		return false, err
	}
	for _, p := range ports {
		if p == port {
			return true, nil
		}
	}
	return false, nil
}

func parsePort(s string) (int, error) {
	s = strings.TrimSpace(s)
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("not a number")
	}
	if n < minPort || n > maxPort {
		return 0, fmt.Errorf("out of range [%d, %d]", minPort, maxPort)
	}
	return n, nil
}

func sorted(m map[int]struct{}) []int {
	out := make([]int, 0, len(m))
	for p := range m {
		out = append(out, p)
	}
	// simple insertion sort — port lists are typically small
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}
