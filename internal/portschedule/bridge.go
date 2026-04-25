package portschedule

import "time"

// Guard wraps a scan function and skips execution when the schedule
// does not permit it at the current time.
type Guard struct {
	sched *Schedule
	now   func() time.Time
}

// NewGuard creates a Guard using the provided Schedule.
// If now is nil, time.Now is used.
func NewGuard(s *Schedule, now func() time.Time) *Guard {
	if now == nil {
		now = time.Now
	}
	return &Guard{sched: s, now: now}
}

// ShouldRun reports whether the scan should execute right now.
func (g *Guard) ShouldRun() bool {
	if g.sched == nil {
		return true
	}
	return g.sched.Allowed(g.now())
}

// Run calls fn only when the schedule permits; otherwise it is a no-op.
// It returns true when fn was called.
func (g *Guard) Run(fn func()) bool {
	if !g.ShouldRun() {
		return false
	}
	fn()
	return true
}
