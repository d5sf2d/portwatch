// Package portstreak tracks consecutive scan cycles in which a port
// remains in a given state (open or closed) for a specific host.
package portstreak

import (
	"fmt"
	"sync"
	"time"
)

// State represents whether a port is open or closed.
type State int

const (
	Open   State = 1
	Closed State = 0
)

// entry holds the current streak information for a single host+port pair.
type entry struct {
	state   State
	count   int
	first   time.Time
	last    time.Time
}

// Streak holds streak data for all tracked host+port pairs.
type Streak struct {
	mu      sync.Mutex
	entries map[string]*entry
	now     func() time.Time
}

// New returns a new Streak tracker.
func New() *Streak {
	return &Streak{
		entries: make(map[string]*entry),
		now:     time.Now,
	}
}

func entryKey(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

// Record records the current state of a port for a host.
// If the state matches the previous state the streak count increments;
// otherwise the streak resets to 1.
func (s *Streak) Record(host string, port int, state State) {
	s.mu.Lock()
	defer s.mu.Unlock()

	k := entryKey(host, port)
	now := s.now()

	if e, ok := s.entries[k]; ok && e.state == state {
		e.count++
		e.last = now
		return
	}

	s.entries[k] = &entry{
		state: state,
		count: 1,
		first: now,
		last:  now,
	}
}

// Count returns the current streak length for a host+port pair.
// Returns 0 if the pair has never been recorded.
func (s *Streak) Count(host string, port int) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	if e, ok := s.entries[entryKey(host, port)]; ok {
		return e.count
	}
	return 0
}

// Current returns the current tracked state and streak count.
// ok is false if the pair has never been recorded.
func (s *Streak) Current(host string, port int) (state State, count int, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, found := s.entries[entryKey(host, port)]
	if !found {
		return 0, 0, false
	}
	return e.state, e.count, true
}

// Reset clears the streak for a specific host+port pair.
func (s *Streak) Reset(host string, port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, entryKey(host, port))
}
