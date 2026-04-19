package scanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// PortState represents the state of a scanned port.
type PortState struct {
	Port     int
	Protocol string
	Open     bool
}

// ScanResult holds the results of a full port scan.
type ScanResult struct {
	Host      string
	ScannedAt time.Time
	Ports     []PortState
}

// Scanner performs TCP port scans on a target host.
type Scanner struct {
	Timeout     time.Duration
	Concurrency int
}

// New creates a Scanner with sensible defaults.
func New() *Scanner {
	return &Scanner{
		Timeout:     500 * time.Millisecond,
		Concurrency: 100,
	}
}

// Scan checks the given ports on host and returns a ScanResult.
func (s *Scanner) Scan(host string, ports []int) (*ScanResult, error) {
	result := &ScanResult{
		Host:      host,
		ScannedAt: time.Now(),
	}

	type job struct{ port int }

	jobs := make(chan job, len(ports))
	results := make(chan PortState, len(ports))

	var wg sync.WaitGroup
	for i := 0; i < s.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				addr := fmt.Sprintf("%s:%d", host, j.port)
				conn, err := net.DialTimeout("tcp", addr, s.Timeout)
				open := err == nil
				if open {
					conn.Close()
				}
				results <- PortState{Port: j.port, Protocol: "tcp", Open: open}
			}
		}()
	}

	for _, p := range ports {
		jobs <- job{port: p}
	}
	close(jobs)

	wg.Wait()
	close(results)

	for ps := range results {
		result.Ports = append(result.Ports, ps)
	}

	return result, nil
}
