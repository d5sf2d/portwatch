// Package ratelimit provides a simple token-bucket rate limiter
// to prevent alert flooding when many ports change simultaneously.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter controls how many events can pass through per interval.
type Limiter struct {
	mu       sync.Mutex
	max      int
	interval time.Duration
	tokens   int
	lastReset time.Time
}

// New creates a Limiter that allows at most max events per interval.
func New(max int, interval time.Duration) *Limiter {
	return &Limiter{
		max:      max,
		interval: interval,
		tokens:   max,
		lastReset: time.Now(),
	}
}

// Allow returns true if the event is permitted under the rate limit.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if now.Sub(l.lastReset) >= l.interval {
		l.tokens = l.max
		l.lastReset = now
	}

	if l.tokens <= 0 {
		return false
	}
	l.tokens--
	return true
}

// Remaining returns the number of tokens left in the current window.
func (l *Limiter) Remaining() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	if time.Since(l.lastReset) >= l.interval {
		return l.max
	}
	return l.tokens
}

// Reset forces the limiter back to a full token bucket immediately.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.tokens = l.max
	l.lastReset = time.Now()
}
