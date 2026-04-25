package portschedule

import (
	"testing"
	"time"
)

func at(weekday time.Weekday, hour, minute int) time.Time {
	// Find a date that matches the requested weekday.
	base := time.Date(2024, 1, 7, hour, minute, 0, 0, time.UTC) // 2024-01-07 is a Sunday
	offset := int(weekday) - int(base.Weekday())
	return base.AddDate(0, 0, offset)
}

func TestAllowed_WithinWindow(t *testing.T) {
	s, err := New([]Window{{Start: 9 * time.Hour, End: 17 * time.Hour}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !s.Allowed(at(time.Monday, 10, 0)) {
		t.Error("expected allowed at 10:00 Mon")
	}
}

func TestAllowed_OutsideWindow(t *testing.T) {
	s, _ := New([]Window{{Start: 9 * time.Hour, End: 17 * time.Hour}}, nil)
	if s.Allowed(at(time.Monday, 18, 0)) {
		t.Error("expected blocked at 18:00")
	}
}

func TestAllowed_WrongWeekday(t *testing.T) {
	s, _ := New([]Window{{Start: 9 * time.Hour, End: 17 * time.Hour}},
		[]time.Weekday{time.Monday, time.Tuesday})
	if s.Allowed(at(time.Saturday, 10, 0)) {
		t.Error("expected blocked on Saturday")
	}
}

func TestAllowed_NoWindows_AnyTime(t *testing.T) {
	s, _ := New(nil, nil)
	if !s.Allowed(at(time.Wednesday, 3, 0)) {
		t.Error("expected allowed when no windows defined")
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New([]Window{{Start: 17 * time.Hour, End: 9 * time.Hour}}, nil)
	if err == nil {
		t.Error("expected error for inverted window")
	}
}

func TestNextAllowed_FindsNextWindow(t *testing.T) {
	s, _ := New([]Window{{Start: 9 * time.Hour, End: 10 * time.Hour}}, nil)
	now := at(time.Monday, 11, 0) // just after window
	next := s.NextAllowed(now)
	if next.IsZero() {
		t.Fatal("expected non-zero next time")
	}
	if next.Hour() != 9 {
		t.Errorf("expected next at 09:xx, got %02d:%02d", next.Hour(), next.Minute())
	}
}

func TestGuard_ShouldRun_Allowed(t *testing.T) {
	s, _ := New([]Window{{Start: 8 * time.Hour, End: 18 * time.Hour}}, nil)
	now := func() time.Time { return at(time.Wednesday, 12, 0) }
	g := NewGuard(s, now)
	if !g.ShouldRun() {
		t.Error("expected ShouldRun true")
	}
}

func TestGuard_Run_SkipsWhenBlocked(t *testing.T) {
	s, _ := New([]Window{{Start: 8 * time.Hour, End: 9 * time.Hour}}, nil)
	now := func() time.Time { return at(time.Wednesday, 22, 0) }
	g := NewGuard(s, now)
	called := false
	ran := g.Run(func() { called = true })
	if ran || called {
		t.Error("expected fn to be skipped outside window")
	}
}

func TestGuard_NilSchedule_AlwaysRuns(t *testing.T) {
	g := NewGuard(nil, nil)
	if !g.ShouldRun() {
		t.Error("expected nil schedule to always allow")
	}
}
