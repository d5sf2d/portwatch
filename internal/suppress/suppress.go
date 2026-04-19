// Package suppress provides a mechanism to silence alerts for specific ports
// during defined time windows (e.g. maintenance windows).
package suppress

import (
	"time"
)

// Window represents a time range during which alerts for a port are suppressed.
type Window struct {
	Port  int
	Start time.Time
	End   time.Time
	Note  string
}

// List holds a collection of suppression windows.
type List struct {
	windows []Window
}

// New returns an empty suppression List.
func New() *List {
	return &List{}
}

// Add registers a suppression window.
func (l *List) Add(w Window) {
	l.windows = append(l.windows, w)
}

// IsSuppressed reports whether alerts for the given port should be silenced at t.
func (l *List) IsSuppressed(port int, t time.Time) bool {
	for _, w := range l.windows {
		if w.Port == port && !t.Before(w.Start) && t.Before(w.End) {
			return true
		}
	}
	return false
}

// Active returns all windows that are currently active at time t.
func (l *List) Active(t time.Time) []Window {
	var out []Window
	for _, w := range l.windows {
		if !t.Before(w.Start) && t.Before(w.End) {
			out = append(out, w)
		}
	}
	return out
}

// Remove deletes all suppression windows for a given port.
func (l *List) Remove(port int) {
	filtered := l.windows[:0]
	for _, w := range l.windows {
		if w.Port != port {
			filtered = append(filtered, w)
		}
	}
	l.windows = filtered
}
