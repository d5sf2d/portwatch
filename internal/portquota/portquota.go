// Package portquota enforces per-host limits on the number of simultaneously
// open ports that portwatch will alert on. Scans that exceed the quota are
// flagged so callers can decide whether to suppress or escalate.
package portquota

import (
	"fmt"
	"sync"
)

// Quota holds the maximum number of open ports allowed for a host.
type Quota struct {
	mu      sync.RWMutex
	limits  map[string]int // host -> max open ports
	default_ int
}

// New returns a Quota with the given default limit applied to any host that
// does not have an explicit override. A default of 0 means unlimited.
func New(defaultLimit int) *Quota {
	return &Quota{
		limits:   make(map[string]int),
		default_: defaultLimit,
	}
}

// SetHost registers an explicit port limit for host. A limit of 0 removes any
// override and falls back to the default.
func (q *Quota) SetHost(host string, limit int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if limit <= 0 {
		delete(q.limits, host)
		return
	}
	q.limits[host] = limit
}

// Limit returns the effective limit for host.
func (q *Quota) Limit(host string) int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if v, ok := q.limits[host]; ok {
		return v
	}
	return q.default_
}

// Check returns an error when openCount exceeds the quota for host. A limit of
// 0 is treated as unlimited and Check always returns nil.
func (q *Quota) Check(host string, openCount int) error {
	limit := q.Limit(host)
	if limit == 0 {
		return nil
	}
	if openCount > limit {
		return fmt.Errorf("portquota: host %q has %d open ports, limit is %d", host, openCount, limit)
	}
	return nil
}

// Exceeded reports whether openCount is strictly greater than the host limit.
func (q *Quota) Exceeded(host string, openCount int) bool {
	return q.Check(host, openCount) != nil
}
