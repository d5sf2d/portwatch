// Package portping provides round-trip latency measurement for open ports.
package portping

import (
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a single ping attempt against a port.
type Result struct {
	Host    string
	Port    int
	Latency time.Duration
	Alive   bool
	Err     error
}

// Pinger measures TCP round-trip latency for a set of ports on a host.
type Pinger struct {
	timeout time.Duration
	count   int
}

// New returns a Pinger with the given timeout per attempt and number of probes.
func New(timeout time.Duration, count int) *Pinger {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	if count <= 0 {
		count = 3
	}
	return &Pinger{timeout: timeout, count: count}
}

// Ping measures the average TCP connect latency for host:port over p.count attempts.
func (p *Pinger) Ping(host string, port int) Result {
	addr := fmt.Sprintf("%s:%d", host, port)
	var total time.Duration
	successes := 0

	for i := 0; i < p.count; i++ {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", addr, p.timeout)
		elapsed := time.Since(start)
		if err != nil {
			return Result{Host: host, Port: port, Alive: false, Err: err}
		}
		conn.Close()
		total += elapsed
		successes++
	}

	return Result{
		Host:    host,
		Port:    port,
		Latency: total / time.Duration(successes),
		Alive:   true,
	}
}

// PingAll pings every port in the list and returns one Result per port.
func (p *Pinger) PingAll(host string, ports []int) []Result {
	results := make([]Result, 0, len(ports))
	for _, port := range ports {
		results = append(results, p.Ping(host, port))
	}
	return results
}
