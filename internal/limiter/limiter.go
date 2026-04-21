// Package limiter provides a token-bucket style alert limiter that caps
// the number of notifications sent per host within a rolling time window.
// This prevents alert storms when many ports change state simultaneously.
package limiter

import (
	"sync"
	"time"
)

// Limiter tracks per-host alert counts within a rolling window and enforces
// a maximum number of alerts allowed in that window.
type Limiter struct {
	mu      sync.Mutex
	window  time.Duration
	max     int
	entries map[string][]time.Time
	now     func() time.Time
}

// New creates a Limiter that allows at most max alerts per host within window.
// Alerts beyond the cap are silently dropped until the window rolls forward.
func New(max int, window time.Duration) *Limiter {
	return &Limiter{
		window:  window,
		max:     max,
		entries: make(map[string][]time.Time),
		now:     time.Now,
	}
}

// Allow reports whether an alert for host should be forwarded.
// It records the attempt and returns false once the per-host cap is reached
// for the current window.
func (l *Limiter) Allow(host string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	l.prune(host, now)

	if len(l.entries[host]) >= l.max {
		return false
	}

	l.entries[host] = append(l.entries[host], now)
	return true
}

// Remaining returns how many more alerts are permitted for host in the current
// window. A value of zero means the host is currently rate-limited.
func (l *Limiter) Remaining(host string) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.prune(host, l.now())
	rem := l.max - len(l.entries[host])
	if rem < 0 {
		return 0
	}
	return rem
}

// Reset clears the alert history for host, immediately restoring its full
// quota. Useful after a manual acknowledgement or test run.
func (l *Limiter) Reset(host string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, host)
}

// ResetAll clears alert history for every tracked host.
func (l *Limiter) ResetAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = make(map[string][]time.Time)
}

// prune removes timestamps that have fallen outside the rolling window.
// Must be called with l.mu held.
func (l *Limiter) prune(host string, now time.Time) {
	cutoff := now.Add(-l.window)
	times := l.entries[host]
	valid := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	if len(valid) == 0 {
		delete(l.entries, host)
	} else {
		l.entries[host] = valid
	}
}
