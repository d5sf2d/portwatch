// Package portmute provides a mute list for suppressing alerts on specific
// ports for a defined duration. Unlike suppress, mute entries are explicitly
// managed by the operator rather than derived from scan behaviour.
package portmute

import (
	"sync"
	"time"
)

// Entry represents a single muted port with an optional expiry.
type Entry struct {
	Port      int
	Reason    string
	ExpiresAt time.Time // zero means indefinite
}

// List holds the current set of muted ports.
type List struct {
	mu      sync.RWMutex
	entries map[int]Entry
	now     func() time.Time
}

// New returns an initialised mute List.
func New() *List {
	return &List{
		entries: make(map[int]Entry),
		now:     time.Now,
	}
}

// Mute adds a port to the mute list for the given duration.
// Pass 0 to mute indefinitely.
func (l *List) Mute(port int, reason string, duration time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	var exp time.Time
	if duration > 0 {
		exp = l.now().Add(duration)
	}
	l.entries[port] = Entry{Port: port, Reason: reason, ExpiresAt: exp}
}

// Unmute removes a port from the mute list immediately.
func (l *List) Unmute(port int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, port)
}

// IsMuted reports whether the given port is currently muted.
func (l *List) IsMuted(port int) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	e, ok := l.entries[port]
	if !ok {
		return false
	}
	if !e.ExpiresAt.IsZero() && l.now().After(e.ExpiresAt) {
		return false
	}
	return true
}

// Active returns all currently active mute entries, pruning expired ones.
func (l *List) Active() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	var out []Entry
	for port, e := range l.entries {
		if !e.ExpiresAt.IsZero() && now.After(e.ExpiresAt) {
			delete(l.entries, port)
			continue
		}
		out = append(out, e)
	}
	return out
}
