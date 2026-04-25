// Package porttrend tracks port open/close frequency over a sliding window
// to identify ports that are unstable or frequently toggling state.
package porttrend

import (
	"sync"
	"time"
)

// Event records a single port state change.
type Event struct {
	Port      int
	Host      string
	OpenedAt  time.Time
	WasClosed bool // true if the port was closed in this event
}

// Summary holds trend statistics for a single port.
type Summary struct {
	Port       int
	Host       string
	OpenCount  int
	CloseCount int
	LastSeen   time.Time
}

// Tracker maintains a sliding-window log of port events per host+port.
type Tracker struct {
	mu     sync.Mutex
	window time.Duration
	events []Event
}

// New returns a Tracker that retains events within the given window.
func New(window time.Duration) *Tracker {
	return &Tracker{window: window}
}

// Record appends an event and prunes entries outside the window.
func (t *Tracker) Record(e Event) {
	if e.OpenedAt.IsZero() {
		e.OpenedAt = time.Now()
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = append(t.events, e)
	t.prune(time.Now())
}

// Trends returns a summary per unique host+port pair seen within the window.
func (t *Tracker) Trends(now time.Time) []Summary {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.prune(now)

	type key struct{ host string; port int }
	m := make(map[key]*Summary)
	for _, e := range t.events {
		k := key{e.Host, e.Port}
		s, ok := m[k]
		if !ok {
			s = &Summary{Port: e.Port, Host: e.Host}
			m[k] = s
		}
		if e.WasClosed {
			s.CloseCount++
		} else {
			s.OpenCount++
		}
		if e.OpenedAt.After(s.LastSeen) {
			s.LastSeen = e.OpenedAt
		}
	}
	out := make([]Summary, 0, len(m))
	for _, s := range m {
		out = append(out, *s)
	}
	return out
}

// prune removes events older than the window. Caller must hold mu.
func (t *Tracker) prune(now time.Time) {
	cutoff := now.Add(-t.window)
	i := 0
	for i < len(t.events) && t.events[i].OpenedAt.Before(cutoff) {
		i++
	}
	t.events = t.events[i:]
}
