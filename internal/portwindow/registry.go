package portwindow

import (
	"fmt"
	"sync"
	"time"
)

// Registry manages one Window per host:port pair.
type Registry struct {
	mu       sync.Mutex
	windows  map[string]*Window
	duration time.Duration
}

// NewRegistry returns a Registry where each Window spans d.
func NewRegistry(d time.Duration) *Registry {
	return &Registry{
		windows:  make(map[string]*Window),
		duration: d,
	}
}

func (r *Registry) key(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

// Record records an open or close event for the given host and port.
func (r *Registry) Record(host string, port int, opened bool) {
	k := r.key(host, port)
	r.mu.Lock()
	w, ok := r.windows[k]
	if !ok {
		w = New(r.duration)
		r.windows[k] = w
	}
	r.mu.Unlock()
	w.Record(opened)
}

// Get returns the Window for the given host and port, or nil if none exists.
func (r *Registry) Get(host string, port int) *Window {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.windows[r.key(host, port)]
}

// ActiveKeys returns all host:port keys that have at least one event in window.
func (r *Registry) ActiveKeys() []string {
	r.mu.Lock()
	keys := make([]string, 0, len(r.windows))
	for k := range r.windows {
		keys = append(keys, k)
	}
	r.mu.Unlock()

	active := keys[:0]
	for _, k := range keys {
		r.mu.Lock()
		w := r.windows[k]
		r.mu.Unlock()
		if w != nil && w.Total() > 0 {
			active = append(active, k)
		}
	}
	return active
}
