// Package portcap tracks the capacity (maximum allowed open ports) per host
// and reports whether a host is over its configured cap.
package portcap

import (
	"errors"
	"fmt"
	"sync"
)

// ErrInvalidCap is returned when a cap value is not positive.
var ErrInvalidCap = errors.New("portcap: cap must be greater than zero")

// Tracker holds per-host port capacity limits.
type Tracker struct {
	mu         sync.RWMutex
	defaultCap int
	hosts      map[string]int
}

// New creates a Tracker with the given default cap applied to all hosts
// unless overridden via SetHost.
func New(defaultCap int) (*Tracker, error) {
	if defaultCap <= 0 {
		return nil, ErrInvalidCap
	}
	return &Tracker{
		defaultCap: defaultCap,
		hosts:      make(map[string]int),
	}, nil
}

// SetHost overrides the cap for a specific host.
func (t *Tracker) SetHost(host string, cap int) error {
	if cap <= 0 {
		return fmt.Errorf("%w: got %d for host %q", ErrInvalidCap, cap, host)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.hosts[host] = cap
	return nil
}

// Cap returns the effective cap for the given host.
func (t *Tracker) Cap(host string) int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if c, ok := t.hosts[host]; ok {
		return c
	}
	return t.defaultCap
}

// Exceeded reports whether the number of open ports surpasses the cap for host.
func (t *Tracker) Exceeded(host string, openCount int) bool {
	return openCount > t.Cap(host)
}

// Overage returns how many ports exceed the cap (0 if within limit).
func (t *Tracker) Overage(host string, openCount int) int {
	over := openCount - t.Cap(host)
	if over < 0 {
		return 0
	}
	return over
}
