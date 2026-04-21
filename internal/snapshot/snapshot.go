// Package snapshot provides utilities for capturing and comparing
// point-in-time views of open ports across one or more hosts.
package snapshot

import (
	"fmt"
	"sort"
	"time"
)

// Port represents a single open port observed during a scan.
type Port struct {
	Number   int    `json:"number"`
	Protocol string `json:"protocol"` // "tcp" or "udp"
	Service  string `json:"service,omitempty"`
}

// Snapshot is a point-in-time capture of open ports on a host.
type Snapshot struct {
	Host      string    `json:"host"`
	Ports     []Port    `json:"ports"`
	CapturedAt time.Time `json:"captured_at"`
}

// New creates a new Snapshot for the given host and ports, stamped with now.
func New(host string, ports []Port) Snapshot {
	sorted := make([]Port, len(ports))
	copy(sorted, ports)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Number != sorted[j].Number {
			return sorted[i].Number < sorted[j].Number
		}
		return sorted[i].Protocol < sorted[j].Protocol
	})
	return Snapshot{
		Host:       host,
		Ports:      sorted,
		CapturedAt: time.Now().UTC(),
	}
}

// PortSet returns the ports as a map keyed by "proto:number" for fast lookup.
func (s Snapshot) PortSet() map[string]Port {
	m := make(map[string]Port, len(s.Ports))
	for _, p := range s.Ports {
		m[key(p)] = p
	}
	return m
}

// Summary returns a human-readable one-line description of the snapshot.
func (s Snapshot) Summary() string {
	return fmt.Sprintf("host=%s ports=%d captured_at=%s",
		s.Host, len(s.Ports), s.CapturedAt.Format(time.RFC3339))
}

func key(p Port) string {
	return fmt.Sprintf("%s:%d", p.Protocol, p.Number)
}
