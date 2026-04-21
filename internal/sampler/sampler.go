// Package sampler provides adaptive scan interval adjustment based on
// observed change frequency. When ports change often the interval shrinks;
// when the host is stable it grows back toward the configured maximum.
package sampler

import (
	"time"
)

const (
	defaultMinInterval = 5 * time.Second
	defaultMaxInterval = 5 * time.Minute
	defaultStepDown    = 0.5  // multiply current by this on change
	defaultStepUp      = 1.25 // multiply current by this on stability
)

// Sampler adjusts a polling interval dynamically.
type Sampler struct {
	min     time.Duration
	max     time.Duration
	current time.Duration
	stepDown float64
	stepUp   float64
}

// Options configures a Sampler.
type Options struct {
	Min      time.Duration
	Max      time.Duration
	Initial  time.Duration
	StepDown float64 // fraction to multiply by on change (0 < v < 1)
	StepUp   float64 // fraction to multiply by on stability (v > 1)
}

// New creates a Sampler with the given options. Zero values fall back to
// sensible defaults.
func New(o Options) *Sampler {
	if o.Min <= 0 {
		o.Min = defaultMinInterval
	}
	if o.Max <= 0 {
		o.Max = defaultMaxInterval
	}
	if o.Initial <= 0 {
		o.Initial = o.Max
	}
	if o.StepDown <= 0 || o.StepDown >= 1 {
		o.StepDown = defaultStepDown
	}
	if o.StepUp <= 1 {
		o.StepUp = defaultStepUp
	}
	current := o.Initial
	if current < o.Min {
		current = o.Min
	}
	if current > o.Max {
		current = o.Max
	}
	return &Sampler{
		min:     o.Min,
		max:     o.Max,
		current: current,
		stepDown: o.StepDown,
		stepUp:   o.StepUp,
	}
}

// Current returns the current recommended polling interval.
func (s *Sampler) Current() time.Duration { return s.current }

// RecordChange notifies the sampler that a port change was detected;
// the interval is reduced toward the minimum.
func (s *Sampler) RecordChange() {
	next := time.Duration(float64(s.current) * s.stepDown)
	if next < s.min {
		next = s.min
	}
	s.current = next
}

// RecordStable notifies the sampler that no change was detected;
// the interval grows back toward the maximum.
func (s *Sampler) RecordStable() {
	next := time.Duration(float64(s.current) * s.stepUp)
	if next > s.max {
		next = s.max
	}
	s.current = next
}

// Reset restores the interval to the configured maximum.
func (s *Sampler) Reset() { s.current = s.max }

// Progress returns a value in [0.0, 1.0] representing how far the current
// interval is between the minimum and maximum. A value of 0.0 means the
// sampler is at its most aggressive (min interval); 1.0 means it is fully
// relaxed (max interval).
func (s *Sampler) Progress() float64 {
	span := float64(s.max - s.min)
	if span == 0 {
		return 1.0
	}
	return float64(s.current-s.min) / span
}
