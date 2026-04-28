// Package portretry provides retry logic for transient port scan failures.
// It wraps a scan attempt and retries up to a configured maximum, backing
// off between attempts, before reporting a definitive open/closed result.
package portretry

import (
	"time"
)

// Result holds the outcome of a retried probe.
type Result struct {
	Port    int
	Open    bool
	Attempts int
}

// Prober is the function signature used to test whether a port is open.
type Prober func(port int) bool

// Retryer retries a port probe up to MaxAttempts times.
type Retryer struct {
	MaxAttempts int
	Backoff     time.Duration
	clock       func(time.Duration)
}

// New returns a Retryer with the given maximum attempts and backoff duration.
// maxAttempts must be >= 1; values below 1 are clamped to 1.
func New(maxAttempts int, backoff time.Duration) *Retryer {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	return &Retryer{
		MaxAttempts: maxAttempts,
		Backoff:     backoff,
		clock:       time.Sleep,
	}
}

// Probe runs the provided prober for port, retrying on a closed result up to
// MaxAttempts times. A port is considered open if any attempt succeeds.
func (r *Retryer) Probe(port int, prober Prober) Result {
	for attempt := 1; attempt <= r.MaxAttempts; attempt++ {
		if prober(port) {
			return Result{Port: port, Open: true, Attempts: attempt}
		}
		if attempt < r.MaxAttempts {
			r.clock(r.Backoff)
		}
	}
	return Result{Port: port, Open: false, Attempts: r.MaxAttempts}
}

// ProbeAll runs Probe for each port in ports and returns all results.
func (r *Retryer) ProbeAll(ports []int, prober Prober) []Result {
	results := make([]Result, 0, len(ports))
	for _, p := range ports {
		results = append(results, r.Probe(p, prober))
	}
	return results
}
