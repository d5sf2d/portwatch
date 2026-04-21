package cooldown

import (
	"testing"
	"time"
)

func newFrozen(base time.Time) (*Tracker, *time.Time) {
	t := New(5 * time.Second)
	current := base
	t.nowFn = func() time.Time { return current }
	return t, &current
}

func TestAllow_FirstCallAllowed(t *testing.T) {
	tr, _ := newFrozen(time.Now())
	if !tr.Allow(8080) {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_WithinWindowBlocked(t *testing.T) {
	base := time.Now()
	tr, current := newFrozen(base)
	tr.Allow(8080)
	*current = base.Add(2 * time.Second)
	if tr.Allow(8080) {
		t.Fatal("expected call within window to be blocked")
	}
}

func TestAllow_AfterWindowAllowed(t *testing.T) {
	base := time.Now()
	tr, current := newFrozen(base)
	tr.Allow(8080)
	*current = base.Add(6 * time.Second)
	if !tr.Allow(8080) {
		t.Fatal("expected call after window to be allowed")
	}
}

func TestAllow_DifferentPortsIndependent(t *testing.T) {
	base := time.Now()
	tr, _ := newFrozen(base)
	tr.Allow(8080)
	if !tr.Allow(9090) {
		t.Fatal("different port should not be affected by cooldown on 8080")
	}
}

func TestReset_ClearsPort(t *testing.T) {
	base := time.Now()
	tr, current := newFrozen(base)
	tr.Allow(8080)
	*current = base.Add(1 * time.Second)
	tr.Reset(8080)
	if !tr.Allow(8080) {
		t.Fatal("expected port to be allowed after reset")
	}
}

func TestActive_ReturnsPortsInCooldown(t *testing.T) {
	base := time.Now()
	tr, current := newFrozen(base)
	tr.Allow(8080)
	tr.Allow(443)
	*current = base.Add(2 * time.Second)
	active := tr.Active()
	if len(active) != 2 {
		t.Fatalf("expected 2 active entries, got %d", len(active))
	}
}

func TestPurge_RemovesExpiredEntries(t *testing.T) {
	base := time.Now()
	tr, current := newFrozen(base)
	tr.Allow(8080)
	*current = base.Add(10 * time.Second)
	tr.Purge()
	if len(tr.entries) != 0 {
		t.Fatalf("expected entries to be purged, got %d", len(tr.entries))
	}
}
