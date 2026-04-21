// Package rollup aggregates multiple port-change diffs within a time
// window into a single batched summary, reducing notification noise when
// many ports change at once.
package rollup

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Rollup accumulates diffs and flushes them after a window elapses or
// the buffer reaches its capacity.
type Rollup struct {
	mu       sync.Mutex
	window   time.Duration
	cap      int
	buf      []state.Diff
	deadline time.Time
	now      func() time.Time
}

// New creates a Rollup with the given flush window and max buffer capacity.
// If cap <= 0 it defaults to 50.
func New(window time.Duration, cap int) *Rollup {
	if cap <= 0 {
		cap = 50
	}
	return &Rollup{
		window: window,
		cap:    cap,
		now:    time.Now,
	}
}

// Add appends a diff to the internal buffer. It returns true and the
// accumulated batch when the window has elapsed or the buffer is full;
// otherwise it returns false and nil.
func (r *Rollup) Add(d state.Diff) ([]state.Diff, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.now()
	if len(r.buf) == 0 {
		r.deadline = now.Add(r.window)
	}

	r.buf = append(r.buf, d)

	if len(r.buf) >= r.cap || now.After(r.deadline) {
		return r.flush(), true
	}
	return nil, false
}

// Flush drains the buffer unconditionally and returns whatever was
// accumulated. Returns nil if the buffer is empty.
func (r *Rollup) Flush() []state.Diff {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.flush()
}

// Pending returns the number of diffs currently buffered.
func (r *Rollup) Pending() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.buf)
}

func (r *Rollup) flush() []state.Diff {
	if len(r.buf) == 0 {
		return nil
	}
	out := make([]state.Diff, len(r.buf))
	copy(out, r.buf)
	r.buf = r.buf[:0]
	return out
}
