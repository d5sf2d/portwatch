// Package portblock provides a registry of explicitly blocked ports that
// should never trigger alerts or appear in scan results.
package portblock

import (
	"fmt"
	"sync"
)

// Registry holds a set of blocked ports with optional reasons.
type Registry struct {
	mu      sync.RWMutex
	blocked map[int]string
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{
		blocked: make(map[int]string),
	}
}

// Block adds port to the blocked set with an optional reason.
func (r *Registry) Block(port int, reason string) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portblock: invalid port %d", port)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.blocked[port] = reason
	return nil
}

// Unblock removes port from the blocked set. It is a no-op if the port
// was not blocked.
func (r *Registry) Unblock(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.blocked, port)
}

// IsBlocked reports whether port is currently blocked.
func (r *Registry) IsBlocked(port int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.blocked[port]
	return ok
}

// Reason returns the reason a port was blocked, or an empty string if the
// port is not in the registry.
func (r *Registry) Reason(port int) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.blocked[port]
}

// All returns a copy of the current blocked-port map.
func (r *Registry) All() map[int]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[int]string, len(r.blocked))
	for k, v := range r.blocked {
		out[k] = v
	}
	return out
}
