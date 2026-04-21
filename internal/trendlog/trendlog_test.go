package trendlog

import (
	"testing"
	"time"
)

func frozenLog(t *testing.T, window time.Duration, at time.Time) *Log {
	t.Helper()
	l := New(window)
	l.now = func() time.Time { return at }
	return l
}

func TestRecord_AppearsInTrends(t *testing.T) {
	now := time.Now()
	l := frozenLog(t, 5*time.Minute, now)

	l.Record(Entry{Port: 80, Host: "localhost", Opened: true})

	trends := l.Trends()
	if len(trends) != 1 {
		t.Fatalf("expected 1 trend, got %d", len(trends))
	}
	if trends[0].Port != 80 || trends[0].OpenCount != 1 || trends[0].TotalEvents != 1 {
		t.Errorf("unexpected trend: %+v", trends[0])
	}
}

func TestRecord_CountsOpenAndClose(t *testing.T) {
	now := time.Now()
	l := frozenLog(t, 5*time.Minute, now)

	l.Record(Entry{Port: 443, Host: "host1", Opened: true})
	l.Record(Entry{Port: 443, Host: "host1", Opened: false})
	l.Record(Entry{Port: 443, Host: "host1", Opened: true})

	trends := l.Trends()
	if len(trends) != 1 {
		t.Fatalf("expected 1 trend, got %d", len(trends))
	}
	tr := trends[0]
	if tr.OpenCount != 2 || tr.CloseCount != 1 || tr.TotalEvents != 3 {
		t.Errorf("unexpected counts: %+v", tr)
	}
}

func TestRecord_OldEventsExpire(t *testing.T) {
	base := time.Now()
	l := New(2 * time.Minute)

	// Record an event 3 minutes in the past.
	old := base.Add(-3 * time.Minute)
	l.entries = append(l.entries, Entry{Port: 22, Host: "h", ChangedAt: old, Opened: true})

	// Advance clock so the old event is outside the window.
	l.now = func() time.Time { return base }

	trends := l.Trends()
	if len(trends) != 0 {
		t.Errorf("expected expired events to be pruned, got %d trends", len(trends))
	}
}

func TestTrends_MultipleHosts(t *testing.T) {
	now := time.Now()
	l := frozenLog(t, 5*time.Minute, now)

	l.Record(Entry{Port: 80, Host: "host-a", Opened: true})
	l.Record(Entry{Port: 80, Host: "host-b", Opened: true})

	trends := l.Trends()
	if len(trends) != 2 {
		t.Fatalf("expected 2 trends (one per host), got %d", len(trends))
	}
}

func TestRecord_SetsTimestampIfZero(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	l := frozenLog(t, 10*time.Minute, now)

	l.Record(Entry{Port: 8080, Host: "h"}) // zero ChangedAt

	if len(l.entries) != 1 {
		t.Fatal("entry not recorded")
	}
	if !l.entries[0].ChangedAt.Equal(now) {
		t.Errorf("expected timestamp %v, got %v", now, l.entries[0].ChangedAt)
	}
}
