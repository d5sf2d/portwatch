// Package portdelta tracks the rate of change (delta) in open ports
// across successive scans, helping identify hosts with unusually high
// port churn.
package portdelta

import (
	"sync"
	"time"
)

// Entry records the number of port changes observed at a point in time.
type Entry struct {
	Host      string
	Delta     int
	RecordedAt time.Time
}

// Tracker accumulates port change deltas per host within a sliding window.
type Tracker struct {
	mu      sync.Mutex
	window  time.Duration
	entries []Entry
	now     func() time.Time
}

// New returns a Tracker that retains entries within the given window.
func New(window time.Duration) *Tracker {
	return &Tracker{
		window: window,
		now:    time.Now,
	}
}

// Record adds a delta observation for the given host.
func (t *Tracker) Record(host string, delta int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.prune()
	t.entries = append(t.entries, Entry{
		Host:       host,
		Delta:      delta,
		RecordedAt: t.now(),
	})
}

// Total returns the sum of all deltas recorded for host within the window.
func (t *Tracker) Total(host string) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.prune()
	sum := 0
	for _, e := range t.entries {
		if e.Host == host {
			sum += e.Delta
		}
	}
	return sum
}

// Hosts returns all hosts that have at least one entry in the window.
func (t *Tracker) Hosts() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.prune()
	seen := make(map[string]struct{})
	var hosts []string
	for _, e := range t.entries {
		if _, ok := seen[e.Host]; !ok {
			seen[e.Host] = struct{}{}
			hosts = append(hosts, e.Host)
		}
	}
	return hosts
}

// prune removes entries older than the window. Must be called with mu held.
func (t *Tracker) prune() {
	cutoff := t.now().Add(-t.window)
	i := 0
	for i < len(t.entries) && t.entries[i].RecordedAt.Before(cutoff) {
		i++
	}
	t.entries = t.entries[i:]
}
