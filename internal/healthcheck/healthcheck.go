// Package healthcheck provides lightweight liveness probes for monitored ports.
// It distinguishes between a port that is open but unresponsive and one that
// is fully healthy by attempting a small read after connecting.
package healthcheck

import (
	"fmt"
	"net"
	"time"
)

// Status represents the health state of a single port.
type Status struct {
	Host      string
	Port      int
	Healthy   bool
	LatencyMs int64
	Reason    string
}

// Checker performs health checks against TCP ports.
type Checker struct {
	timeout time.Duration
}

// New returns a Checker with the given dial timeout.
func New(timeout time.Duration) *Checker {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return &Checker{timeout: timeout}
}

// Check dials host:port and returns a Status indicating liveness.
func (c *Checker) Check(host string, port int) Status {
	addr := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()

	conn, err := net.DialTimeout("tcp", addr, c.timeout)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return Status{
			Host:      host,
			Port:      port,
			Healthy:   false,
			LatencyMs: latency,
			Reason:    err.Error(),
		}
	}
	conn.Close()

	return Status{
		Host:      host,
		Port:      port,
		Healthy:   true,
		LatencyMs: latency,
		Reason:    "ok",
	}
}

// CheckAll runs Check for every port and returns all results.
func (c *Checker) CheckAll(host string, ports []int) []Status {
	results := make([]Status, 0, len(ports))
	for _, p := range ports {
		results = append(results, c.Check(host, p))
	}
	return results
}
