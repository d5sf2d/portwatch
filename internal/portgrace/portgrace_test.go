package portgrace

import (
	"testing"
	"time"
)

func frozenRegistry(dur time.Duration) (*Registry, time.Time) {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	return New(dur), now
}

func TestObserve_FirstCallRegisters(t *testing.T) {
	r, now := frozenRegistry(30 * time.Second)
	r.Observe("localhost", 8080, now)
	if !r.InGrace("localhost", 8080, now) {
		t.Fatal("expected port to be in grace immediately after observe")
	}
}

func TestObserve_SecondCallIsNoop(t *testing.T) {
	r, now := frozenRegistry(30 * time.Second)
	r.Observe("localhost", 8080, now)
	// second observe at a later time should not reset the window
	later := now.Add(20 * time.Second)
	r.Observe("localhost", 8080, later)
	// window should still expire at now+30s, so at now+35s it must be gone
	expired := now.Add(35 * time.Second)
	if r.InGrace("localhost", 8080, expired) {
		t.Fatal("expected grace to have expired; second observe must not reset it")
	}
}

func TestInGrace_WithinWindow(t *testing.T) {
	r, now := frozenRegistry(1 * time.Minute)
	r.Observe("host", 22, now)
	if !r.InGrace("host", 22, now.Add(30*time.Second)) {
		t.Fatal("expected port to still be in grace at 30s")
	}
}

func TestInGrace_AfterWindow(t *testing.T) {
	r, now := frozenRegistry(10 * time.Second)
	r.Observe("host", 22, now)
	if r.InGrace("host", 22, now.Add(11*time.Second)) {
		t.Fatal("expected grace to have expired after window")
	}
}

func TestInGrace_UnknownPort(t *testing.T) {
	r, now := frozenRegistry(1 * time.Minute)
	if r.InGrace("host", 9999, now) {
		t.Fatal("unknown port should not be in grace")
	}
}

func TestForget_RemovesPort(t *testing.T) {
	r, now := frozenRegistry(1 * time.Minute)
	r.Observe("host", 443, now)
	r.Forget("host", 443)
	if r.InGrace("host", 443, now) {
		t.Fatal("port should not be in grace after Forget")
	}
}

func TestActive_ReturnsOnlyGracePorts(t *testing.T) {
	r, now := frozenRegistry(30 * time.Second)
	r.Observe("host", 80, now)
	r.Observe("host", 443, now)
	r.Observe("host", 22, now.Add(-31*time.Second)) // already expired

	active := r.Active(now)
	if len(active) != 2 {
		t.Fatalf("expected 2 active grace entries, got %d", len(active))
	}
}

func TestActive_EmptyRegistry(t *testing.T) {
	r, now := frozenRegistry(1 * time.Minute)
	if got := r.Active(now); len(got) != 0 {
		t.Fatalf("expected empty active list, got %d entries", len(got))
	}
}
