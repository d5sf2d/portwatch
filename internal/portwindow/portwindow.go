// Package portwindow tracks rolling time-window statistics for port open/close
// events, allowing callers to query how many events occurred in the last N
// seconds for a given host:port pair.
package portwindow

import (
	"sync"
	"time"
)

// event records a single open or close occurrence.
type event struct {
	at     time.Time
	opened bool
}

// Window holds rolling event history for a single host:port key.
type Window struct {
	mu       sync.Mutex
	events   []event
	duration time.Duration
	now      func() time.Time
}

// New creates a Window that retains events within the given duration.
func New(d time.Duration) *Window {
	return &Window{duration: d, now: time.Now}
}

// Record appends an open or close event for the given host:port key.
func (w *Window) Record(opened bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.prune()
	w.events = append(w.events, event{at: w.now(), opened: opened})
}

// Opens returns the count of open events within the window.
func (w *Window) Opens() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.prune()
	count := 0
	for _, e := range w.events {
		if e.opened {
			count++
		}
	}
	return count
}

// Closes returns the count of close events within the window.
func (w *Window) Closes() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.prune()
	count := 0
	for _, e := range w.events {
		if !e.opened {
			count++
		}
	}
	return count
}

// Total returns the total event count within the window.
func (w *Window) Total() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.prune()
	return len(w.events)
}

// prune removes events older than the window duration. Caller must hold w.mu.
func (w *Window) prune() {
	cutoff := w.now().Add(-w.duration)
	i := 0
	for i < len(w.events) && w.events[i].at.Before(cutoff) {
		i++
	}
	w.events = w.events[i:]
}
