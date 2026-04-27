// Package portgrace tracks ports that are in a grace period after first
// appearing, suppressing alerts until the window has elapsed. This avoids
// noise from transient services that open and close quickly on startup.
package portgrace

import (
	"sync"
	"time"
)

// Entry holds the grace-period metadata for a single port on a host.
type Entry struct {
	Host      string
	Port      int
	OpenedAt  time.Time
	Duration  time.Duration
}

// InGrace reports whether the entry is still within its grace window.
func (e Entry) InGrace(now time.Time) bool {
	return now.Before(e.OpenedAt.Add(e.Duration))
}

// Registry tracks grace-period entries keyed by host+port.
type Registry struct {
	mu       sync.Mutex
	entries  map[string]Entry
	default_ time.Duration
}

// New creates a Registry with the given default grace duration.
func New(defaultDuration time.Duration) *Registry {
	return &Registry{
		entries:  make(map[string]Entry),
		default_: defaultDuration,
	}
}

func entryKey(host string, port int) string {
	return host + ":" + itoa(port)
}

// Observe records that a port was seen open at the given time. If the port is
// already tracked the call is a no-op (first-seen wins).
func (r *Registry) Observe(host string, port int, now time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := entryKey(host, port)
	if _, ok := r.entries[k]; ok {
		return
	}
	r.entries[k] = Entry{
		Host:     host,
		Port:     port,
		OpenedAt: now,
		Duration: r.default_,
	}
}

// InGrace reports whether the given port is still within its grace period.
func (r *Registry) InGrace(host string, port int, now time.Time) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.entries[entryKey(host, port)]
	if !ok {
		return false
	}
	return e.InGrace(now)
}

// Forget removes a port from the registry (e.g. when it closes).
func (r *Registry) Forget(host string, port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, entryKey(host, port))
}

// Active returns all entries currently within their grace window.
func (r *Registry) Active(now time.Time) []Entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []Entry
	for _, e := range r.entries {
		if e.InGrace(now) {
			out = append(out, e)
		}
	}
	return out
}

// itoa is a minimal int-to-string helper to avoid importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
