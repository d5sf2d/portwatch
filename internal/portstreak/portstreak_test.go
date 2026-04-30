package portstreak

import (
	"testing"
	"time"
)

func frozenStreak() *Streak {
	s := New()
	t := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	s.now = func() time.Time { return t }
	return s
}

func TestRecord_IncreasesStreakOnSameState(t *testing.T) {
	s := frozenStreak()
	s.Record("localhost", 80, Open)
	s.Record("localhost", 80, Open)
	s.Record("localhost", 80, Open)

	if got := s.Count("localhost", 80); got != 3 {
		t.Fatalf("expected streak 3, got %d", got)
	}
}

func TestRecord_ResetsOnStateChange(t *testing.T) {
	s := frozenStreak()
	s.Record("localhost", 80, Open)
	s.Record("localhost", 80, Open)
	s.Record("localhost", 80, Closed)

	if got := s.Count("localhost", 80); got != 1 {
		t.Fatalf("expected streak 1 after state change, got %d", got)
	}
}

func TestCount_UnknownPortIsZero(t *testing.T) {
	s := frozenStreak()
	if got := s.Count("localhost", 9999); got != 0 {
		t.Fatalf("expected 0 for unknown port, got %d", got)
	}
}

func TestCurrent_ReturnsStateAndCount(t *testing.T) {
	s := frozenStreak()
	s.Record("host-a", 443, Open)
	s.Record("host-a", 443, Open)

	state, count, ok := s.Current("host-a", 443)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if state != Open {
		t.Fatalf("expected Open, got %v", state)
	}
	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func TestCurrent_MissingReturnsFalse(t *testing.T) {
	s := frozenStreak()
	_, _, ok := s.Current("ghost", 22)
	if ok {
		t.Fatal("expected ok=false for unseen host+port")
	}
}

func TestRecord_DifferentPortsAreIndependent(t *testing.T) {
	s := frozenStreak()
	s.Record("localhost", 22, Open)
	s.Record("localhost", 22, Open)
	s.Record("localhost", 80, Open)

	if got := s.Count("localhost", 22); got != 2 {
		t.Fatalf("port 22: expected 2, got %d", got)
	}
	if got := s.Count("localhost", 80); got != 1 {
		t.Fatalf("port 80: expected 1, got %d", got)
	}
}

func TestReset_ClearsEntry(t *testing.T) {
	s := frozenStreak()
	s.Record("localhost", 8080, Open)
	s.Record("localhost", 8080, Open)
	s.Reset("localhost", 8080)

	if got := s.Count("localhost", 8080); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}
