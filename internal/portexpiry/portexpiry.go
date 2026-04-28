// Package portexpiry tracks how long ports have been continuously open
// and flags those that exceed a configured maximum open duration.
package portexpiry

import (
	"sync"
	"time"
)

// Entry holds the first-seen timestamp and the expiry duration for a port.
type Entry struct {
	FirstSeen time.Time
	MaxAge    time.Duration
}

// Tracker monitors port open durations and reports expiry violations.
type Tracker struct {
	mu         sync.Mutex
	entries    map[string]Entry
	defaultMax time.Duration
	now        func() time.Time
}

// New creates a Tracker with the given default maximum open duration.
// A zero defaultMax means ports never expire by default.
func New(defaultMax time.Duration) *Tracker {
	return &Tracker{
		entries:    make(map[string]Entry),
		defaultMax: defaultMax,
		now:        time.Now,
	}
}

func key(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

// Observe records that a port is open. If it has not been seen before,
// the current time is stored as the first-seen timestamp.
func (t *Tracker) Observe(host string, port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	k := key(host, port)
	if _, ok := t.entries[k]; !ok {
		t.entries[k] = Entry{
			FirstSeen: t.now(),
			MaxAge:    t.defaultMax,
		}
	}
}

// Forget removes a port from tracking (e.g. when it closes).
func (t *Tracker) Forget(host string, port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, key(host, port))
}

// SetMaxAge overrides the maximum open duration for a specific port.
func (t *Tracker) SetMaxAge(host string, port int, max time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	k := key(host, port)
	e := t.entries[k]
	e.MaxAge = max
	t.entries[k] = e
}

// Expired reports whether the port has been open longer than its allowed
// maximum. Ports with a zero MaxAge are never considered expired.
func (t *Tracker) Expired(host string, port int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[key(host, port)]
	if !ok || e.MaxAge == 0 {
		return false
	}
	return t.now().Sub(e.FirstSeen) > e.MaxAge
}

// Age returns how long the port has been continuously open.
// Returns 0 if the port is not tracked.
func (t *Tracker) Age(host string, port int) time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[key(host, port)]
	if !ok {
		return 0
	}
	return t.now().Sub(e.FirstSeen)
}
