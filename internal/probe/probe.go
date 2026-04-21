// Package probe performs lightweight TCP reachability checks with
// latency measurement, suitable for feeding into health dashboards.
package probe

import (
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a single probe attempt.
type Result struct {
	Host    string
	Port    int
	Open    bool
	Latency time.Duration
	Err     error
}

// Prober sends TCP probes to host:port combinations.
type Prober struct {
	timeout time.Duration
}

// New returns a Prober with the given dial timeout.
func New(timeout time.Duration) *Prober {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return &Prober{timeout: timeout}
}

// Probe dials host:port and returns a Result.
func (p *Prober) Probe(host string, port int) Result {
	addr := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, p.timeout)
	latency := time.Since(start)
	if err != nil {
		return Result{Host: host, Port: port, Open: false, Latency: latency, Err: err}
	}
	conn.Close()
	return Result{Host: host, Port: port, Open: true, Latency: latency}
}

// ProbeAll probes every port in the list concurrently and returns all results.
func (p *Prober) ProbeAll(host string, ports []int) []Result {
	results := make([]Result, len(ports))
	type indexed struct {
		i int
		r Result
	}
	ch := make(chan indexed, len(ports))
	for i, port := range ports {
		go func(idx, pt int) {
			ch <- indexed{idx, p.Probe(host, pt)}
		}(i, port)
	}
	for range ports {
		item := <-ch
		results[item.i] = item.r
	}
	return results
}
