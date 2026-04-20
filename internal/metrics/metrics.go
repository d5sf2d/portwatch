// Package metrics tracks runtime counters for portwatch scan cycles.
package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of collected counters.
type Snapshot struct {
	ScansTotal      int64
	AlertsTotal     int64
	PortsOpen       int
	SuppressedTotal int64
	LastScanAt      time.Time
	LastScanDur     time.Duration
}

// Collector accumulates scan-cycle metrics in memory.
type Collector struct {
	mu              sync.Mutex
	scansTotal      int64
	alertsTotal     int64
	portsOpen       int
	suppressedTotal int64
	lastScanAt      time.Time
	lastScanDur     time.Duration
}

// New returns a zero-value Collector ready for use.
func New() *Collector {
	return &Collector{}
}

// RecordScan records the completion of one scan cycle.
func (c *Collector) RecordScan(dur time.Duration, openPorts int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.scansTotal++
	c.lastScanAt = time.Now()
	c.lastScanDur = dur
	c.portsOpen = openPorts
}

// RecordAlert increments the alert counter by n.
func (c *Collector) RecordAlert(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.alertsTotal += int64(n)
}

// RecordSuppressed increments the suppressed-alert counter by n.
func (c *Collector) RecordSuppressed(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.suppressedTotal += int64(n)
}

// Snapshot returns a consistent copy of the current counters.
func (c *Collector) Snapshot() Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	return Snapshot{
		ScansTotal:      c.scansTotal,
		AlertsTotal:     c.alertsTotal,
		PortsOpen:       c.portsOpen,
		SuppressedTotal: c.suppressedTotal,
		LastScanAt:      c.lastScanAt,
		LastScanDur:     c.lastScanDur,
	}
}

// Reset zeroes all counters.
func (c *Collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	*c = Collector{}
}
