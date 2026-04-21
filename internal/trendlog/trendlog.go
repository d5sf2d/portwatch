// Package trendlog tracks port open/close frequency over a sliding window
// to identify ports that are changing state unusually often.
package trendlog

import (
	"sync"
	"time"
)

// Entry records a single state-change event for a port.
type Entry struct {
	Port      int
	Host      string
	ChangedAt time.Time
	Opened    bool // true = opened, false = closed
}

// Trend summarises activity for a single port within the window.
type Trend struct {
	Port        int
	Host        string
	OpenCount   int
	CloseCount  int
	TotalEvents int
}

// Log holds recent change events and can compute per-port trends.
type Log struct {
	mu      sync.Mutex
	entries []Entry
	window  time.Duration
	now     func() time.Time
}

// New creates a Log that retains events within the given sliding window.
func New(window time.Duration) *Log {
	return &Log{
		window: window,
		now:    time.Now,
	}
}

// Record appends a new change event, pruning entries outside the window.
func (l *Log) Record(e Entry) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if e.ChangedAt.IsZero() {
		e.ChangedAt = l.now()
	}
	l.entries = append(l.entries, e)
	l.prune()
}

// Trends returns a Trend summary for every port seen within the window.
func (l *Log) Trends() []Trend {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prune()

	type key struct {
		port int
		host string
	}
	m := make(map[key]*Trend)
	for _, e := range l.entries {
		k := key{e.Port, e.Host}
		if _, ok := m[k]; !ok {
			m[k] = &Trend{Port: e.Port, Host: e.Host}
		}
		t := m[k]
		t.TotalEvents++
		if e.Opened {
			t.OpenCount++
		} else {
			t.CloseCount++
		}
	}

	out := make([]Trend, 0, len(m))
	for _, t := range m {
		out = append(out, *t)
	}
	return out
}

// prune removes entries older than the window. Caller must hold l.mu.
func (l *Log) prune() {
	cutoff := l.now().Add(-l.window)
	i := 0
	for i < len(l.entries) && l.entries[i].ChangedAt.Before(cutoff) {
		i++
	}
	l.entries = l.entries[i:]
}
