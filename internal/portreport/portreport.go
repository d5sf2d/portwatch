// Package portreport aggregates per-port statistics across multiple scans
// and produces a structured summary suitable for display or export.
package portreport

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// PortStat holds accumulated statistics for a single port.
type PortStat struct {
	Port       int
	Host       string
	OpenCount  int
	CloseCount int
	LastSeen   time.Time
	FirstSeen  time.Time
	Label      string
}

// Report is a snapshot of all tracked port statistics.
type Report struct {
	GeneratedAt time.Time
	Stats       []PortStat
}

// Tracker accumulates port open/close events and produces reports.
type Tracker struct {
	mu    sync.Mutex
	entries map[string]*PortStat
	now   func() time.Time
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]*PortStat),
		now:     time.Now,
	}
}

func entryKey(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

// RecordOpen records that a port was observed open.
func (t *Tracker) RecordOpen(host string, port int, label string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	k := entryKey(host, port)
	s, ok := t.entries[k]
	if !ok {
		s = &PortStat{Port: port, Host: host, FirstSeen: now, Label: label}
		t.entries[k] = s
	}
	s.OpenCount++
	s.LastSeen = now
	if label != "" {
		s.Label = label
	}
}

// RecordClose records that a port was observed closed.
func (t *Tracker) RecordClose(host string, port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	k := entryKey(host, port)
	s, ok := t.entries[k]
	if !ok {
		s = &PortStat{Port: port, Host: host, FirstSeen: now}
		t.entries[k] = s
	}
	s.CloseCount++
	s.LastSeen = now
}

// Report returns a sorted snapshot of all accumulated statistics.
func (t *Tracker) Report() Report {
	t.mu.Lock()
	defer t.mu.Unlock()
	stats := make([]PortStat, 0, len(t.entries))
	for _, s := range t.entries {
		copy := *s
		stats = append(stats, copy)
	}
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Host != stats[j].Host {
			return stats[i].Host < stats[j].Host
		}
		return stats[i].Port < stats[j].Port
	})
	return Report{GeneratedAt: t.now(), Stats: stats}
}

// Reset clears all accumulated statistics.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[string]*PortStat)
}
