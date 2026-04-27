// Package portage tracks how long ports have been continuously open,
// allowing detection of long-lived or stale services.
package portage

import (
	"sync"
	"time"
)

// Record holds the first-seen timestamp and open duration for a port.
type Record struct {
	Host     string
	Port     int
	FirstSeen time.Time
	Age      time.Duration
}

// Tracker records when each (host, port) pair was first observed open.
type Tracker struct {
	mu    sync.Mutex
	first map[string]time.Time
	now   func() time.Time
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{
		first: make(map[string]time.Time),
		now:   time.Now,
	}
}

func key(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

// Observe marks a port as open. If it was already tracked, the original
// first-seen time is preserved.
func (t *Tracker) Observe(host string, port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	k := key(host, port)
	if _, ok := t.first[k]; !ok {
		t.first[k] = t.now()
	}
}

// Forget removes a port from tracking (e.g. when it closes).
func (t *Tracker) Forget(host string, port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.first, key(host, port))
}

// Age returns the duration a port has been continuously open.
// Returns 0 and false if the port is not tracked.
func (t *Tracker) Age(host string, port int) (time.Duration, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	fs, ok := t.first[key(host, port)]
	if !ok {
		return 0, false
	}
	return t.now().Sub(fs), true
}

// All returns a snapshot of all tracked ports with their ages.
func (t *Tracker) All() []Record {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Record, 0, len(t.first))
	now := t.now()
	for k, fs := range t.first {
		host, port := splitKey(k)
		out = append(out, Record{
			Host:      host,
			Port:      port,
			FirstSeen: fs,
			Age:       now.Sub(fs),
		})
	}
	return out
}
