// Package watchdog provides threshold-based alerting when a port's
// open/closed state flips more than N times within a rolling window.
package watchdog

import (
	"sync"
	"time"
)

// Breach describes a port that has exceeded its flip threshold.
type Breach struct {
	Port      int
	Flips     int
	Window    time.Duration
	Threshold int
}

// Watchdog tracks state-flip counts per port and reports breaches.
type Watchdog struct {
	mu        sync.Mutex
	threshold int
	window    time.Duration
	events    map[int][]time.Time
}

// New returns a Watchdog that fires when a port flips more than threshold
// times inside window.
func New(threshold int, window time.Duration) *Watchdog {
	if threshold <= 0 {
		threshold = 3
	}
	if window <= 0 {
		window = 5 * time.Minute
	}
	return &Watchdog{
		threshold: threshold,
		window:    window,
		events:    make(map[int][]time.Time),
	}
}

// Record registers a state flip for port at the given time and returns a
// Breach if the threshold is exceeded, or nil otherwise.
func (w *Watchdog) Record(port int, at time.Time) *Breach {
	w.mu.Lock()
	defer w.mu.Unlock()

	cutoff := at.Add(-w.window)
	prev := w.events[port]
	filtered := prev[:0]
	for _, t := range prev {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, at)
	w.events[port] = filtered

	if len(filtered) > w.threshold {
		return &Breach{
			Port:      port,
			Flips:     len(filtered),
			Window:    w.window,
			Threshold: w.threshold,
		}
	}
	return nil
}

// Reset clears the flip history for a single port.
func (w *Watchdog) Reset(port int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.events, port)
}

// Active returns a copy of all ports that currently have recorded flips
// within the rolling window, keyed by port number.
func (w *Watchdog) Active(now time.Time) map[int]int {
	w.mu.Lock()
	defer w.mu.Unlock()
	cutoff := now.Add(-w.window)
	out := make(map[int]int)
	for port, times := range w.events {
		count := 0
		for _, t := range times {
			if t.After(cutoff) {
				count++
			}
		}
		if count > 0 {
			out[port] = count
		}
	}
	return out
}
