// Package portschedule provides time-window based scheduling for port scans,
// allowing scans to be restricted to specific hours or days.
package portschedule

import (
	"fmt"
	"time"
)

// Window represents a time range within a day during which scanning is permitted.
type Window struct {
	Start time.Duration // offset from midnight
	End   time.Duration // offset from midnight
}

// Schedule holds a set of permitted scan windows and the days they apply to.
type Schedule struct {
	windows  []Window
	weekdays map[time.Weekday]bool
}

// New creates a Schedule. If no weekdays are provided, all days are permitted.
func New(windows []Window, days []time.Weekday) (*Schedule, error) {
	for _, w := range windows {
		if w.Start >= w.End {
			return nil, fmt.Errorf("portschedule: window start must be before end")
		}
		if w.End > 24*time.Hour {
			return nil, fmt.Errorf("portschedule: window end exceeds 24h")
		}
	}
	weekdays := make(map[time.Weekday]bool)
	if len(days) == 0 {
		for d := time.Sunday; d <= time.Saturday; d++ {
			weekdays[d] = true
		}
	} else {
		for _, d := range days {
			weekdays[d] = true
		}
	}
	return &Schedule{windows: windows, weekdays: weekdays}, nil
}

// Allowed reports whether scanning is permitted at the given time.
func (s *Schedule) Allowed(t time.Time) bool {
	if !s.weekdays[t.Weekday()] {
		return false
	}
	if len(s.windows) == 0 {
		return true
	}
	offset := time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second
	for _, w := range s.windows {
		if offset >= w.Start && offset < w.End {
			return true
		}
	}
	return false
}

// NextAllowed returns the next time at or after t when scanning is permitted.
// It searches up to 8 days ahead before returning the zero time.
func (s *Schedule) NextAllowed(t time.Time) time.Time {
	const maxSearch = 8 * 24 * time.Hour
	step := time.Minute
	for d := time.Duration(0); d < maxSearch; d += step {
		candidate := t.Add(d)
		if s.Allowed(candidate) {
			return candidate
		}
	}
	return time.Time{}
}
