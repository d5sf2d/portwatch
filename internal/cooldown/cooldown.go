// Package cooldown provides per-port cooldown tracking to prevent
// repeated alerts for the same port within a configurable time window.
package cooldown

import (
	"sync"
	"time"
)

// Entry records the last alert time for a port.
type Entry struct {
	Port     int
	LastSeen time.Time
}

// Tracker holds cooldown state for ports.
type Tracker struct {
	mu       sync.Mutex
	window   time.Duration
	entries  map[int]time.Time
	nowFn    func() time.Time
}

// New returns a Tracker with the given cooldown window.
func New(window time.Duration) *Tracker {
	return &Tracker{
		window:  window,
		entries: make(map[int]time.Time),
		nowFn:   time.Now,
	}
}

// Allow returns true if the port is not in cooldown and records the
// current time as the last-seen time for that port.
func (t *Tracker) Allow(port int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.nowFn()
	if last, ok := t.entries[port]; ok {
		if now.Sub(last) < t.window {
			return false
		}
	}
	t.entries[port] = now
	return true
}

// Reset clears the cooldown record for a specific port.
func (t *Tracker) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, port)
}

// Active returns a snapshot of all ports currently in cooldown.
func (t *Tracker) Active() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.nowFn()
	var out []Entry
	for port, last := range t.entries {
		if now.Sub(last) < t.window {
			out = append(out, Entry{Port: port, LastSeen: last})
		}
	}
	return out
}

// Purge removes all expired cooldown entries.
func (t *Tracker) Purge() {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.nowFn()
	for port, last := range t.entries {
		if now.Sub(last) >= t.window {
			delete(t.entries, port)
		}
	}
}
