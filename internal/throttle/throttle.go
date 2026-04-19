// Package throttle limits how frequently alerts are emitted per port.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks the last alert time per port and suppresses duplicates
// within a configurable cooldown window.
type Throttle struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[int]time.Time
}

// New creates a Throttle with the given cooldown duration.
func New(cooldown time.Duration) *Throttle {
	return &Throttle{
		cooldown: cooldown,
		last:     make(map[int]time.Time),
	}
}

// Allow returns true if an alert for the given port should be emitted.
// It updates the last-seen timestamp when it returns true.
func (t *Throttle) Allow(port int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	if ts, ok := t.last[port]; ok && now.Sub(ts) < t.cooldown {
		return false
	}
	t.last[port] = now
	return true
}

// Reset clears the throttle state for a specific port.
func (t *Throttle) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, port)
}

// ResetAll clears all throttle state.
func (t *Throttle) ResetAll() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last = make(map[int]time.Time)
}

// ActivePorts returns the list of ports currently within the cooldown window.
func (t *Throttle) ActivePorts() []int {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	var ports []int
	for p, ts := range t.last {
		if now.Sub(ts) < t.cooldown {
			ports = append(ports, p)
		}
	}
	return ports
}
