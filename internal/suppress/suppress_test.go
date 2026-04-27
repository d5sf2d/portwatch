package suppress

import (
	"testing"
	"time"
)

var (
	now    = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	before = now.Add(-1 * time.Hour)
	after  = now.Add(1 * time.Hour)
)

func makeList() *List {
	l := New()
	l.Add(Window{Port: 8080, Start: before, End: after, Note: "maintenance"})
	return l
}

func TestIsSuppressed_WithinWindow(t *testing.T) {
	l := makeList()
	if !l.IsSuppressed(8080, now) {
		t.Fatal("expected port 8080 to be suppressed")
	}
}

func TestIsSuppressed_OutsideWindow(t *testing.T) {
	l := makeList()
	outside := now.Add(2 * time.Hour)
	if l.IsSuppressed(8080, outside) {
		t.Fatal("expected port 8080 not to be suppressed outside window")
	}
}

func TestIsSuppressed_DifferentPort(t *testing.T) {
	l := makeList()
	if l.IsSuppressed(9090, now) {
		t.Fatal("expected port 9090 not to be suppressed")
	}
}

func TestIsSuppressed_AtWindowBoundaries(t *testing.T) {
	l := makeList()
	// Exactly at Start should be suppressed
	if !l.IsSuppressed(8080, before) {
		t.Fatal("expected port 8080 to be suppressed at window start")
	}
	// Exactly at End should be suppressed
	if !l.IsSuppressed(8080, after) {
		t.Fatal("expected port 8080 to be suppressed at window end")
	}
}

func TestActive_ReturnsCurrentWindows(t *testing.T) {
	l := makeList()
	l.Add(Window{Port: 443, Start: after, End: after.Add(time.Hour)})
	active := l.Active(now)
	if len(active) != 1 {
		t.Fatalf("expected 1 active window, got %d", len(active))
	}
	if active[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", active[0].Port)
	}
}

func TestRemove_ClearsPort(t *testing.T) {
	l := makeList()
	l.Remove(8080)
	if l.IsSuppressed(8080, now) {
		t.Fatal("expected port 8080 to be unsuppressed after removal")
	}
}

func TestRemove_LeavesOtherPorts(t *testing.T) {
	l := makeList()
	l.Add(Window{Port: 443, Start: before, End: after})
	l.Remove(8080)
	if !l.IsSuppressed(443, now) {
		t.Fatal("expected port 443 to remain suppressed")
	}
}
