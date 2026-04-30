package portdelta

import (
	"testing"
	"time"
)

func frozenTracker(ts time.Time) *Tracker {
	t := New(5 * time.Minute)
	t.now = func() time.Time { return ts }
	return t
}

func TestRecord_AppearsInTotal(t *testing.T) {
	base := time.Now()
	tr := frozenTracker(base)
	tr.Record("host-a", 3)
	tr.Record("host-a", 2)
	if got := tr.Total("host-a"); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestRecord_OldEntriesExpire(t *testing.T) {
	base := time.Now()
	tr := New(5 * time.Minute)
	tr.now = func() time.Time { return base }
	tr.Record("host-a", 4)

	// Advance past the window.
	tr.now = func() time.Time { return base.Add(6 * time.Minute) }
	tr.Record("host-a", 1)

	if got := tr.Total("host-a"); got != 1 {
		t.Fatalf("expected 1 after expiry, got %d", got)
	}
}

func TestTotal_UnknownHostIsZero(t *testing.T) {
	tr := frozenTracker(time.Now())
	if got := tr.Total("nobody"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestHosts_ReturnsDistinct(t *testing.T) {
	base := time.Now()
	tr := frozenTracker(base)
	tr.Record("alpha", 1)
	tr.Record("beta", 2)
	tr.Record("alpha", 3)

	hosts := tr.Hosts()
	if len(hosts) != 2 {
		t.Fatalf("expected 2 distinct hosts, got %d", len(hosts))
	}
}

func TestHosts_EmptyWhenAllExpired(t *testing.T) {
	base := time.Now()
	tr := New(1 * time.Minute)
	tr.now = func() time.Time { return base }
	tr.Record("host-x", 5)

	tr.now = func() time.Time { return base.Add(2 * time.Minute) }
	if hosts := tr.Hosts(); len(hosts) != 0 {
		t.Fatalf("expected no hosts after expiry, got %v", hosts)
	}
}

func TestRecord_MultipleHosts_IndependentTotals(t *testing.T) {
	base := time.Now()
	tr := frozenTracker(base)
	tr.Record("h1", 10)
	tr.Record("h2", 7)

	if got := tr.Total("h1"); got != 10 {
		t.Fatalf("h1: expected 10, got %d", got)
	}
	if got := tr.Total("h2"); got != 7 {
		t.Fatalf("h2: expected 7, got %d", got)
	}
}
