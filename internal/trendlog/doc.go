// Package trendlog provides a sliding-window event log that tracks how
// frequently individual ports change state (open ↔ closed) on a host.
//
// Use it to surface "flapping" ports — ports that open and close repeatedly
// within a short period — without relying on the heavier watchdog package.
package trendlog
