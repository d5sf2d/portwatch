package watchdog

import (
	"testing"
	"time"
)

func frozen() time.Time {
	return time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
}

func TestRecord_BelowThreshold(t *testing.T) {
	w := New(3, time.Minute)
	now := frozen()
	for i := 0; i < 3; i++ {
		b := w.Record(80, now.Add(time.Duration(i)*time.Second))
		if b != nil {
			t.Fatalf("expected no breach at flip %d, got %+v", i+1, b)
		}
	}
}

func TestRecord_ExceedsThreshold(t *testing.T) {
	w := New(3, time.Minute)
	now := frozen()
	var breach *Breach
	for i := 0; i < 4; i++ {
		breach = w.Record(443, now.Add(time.Duration(i)*time.Second))
	}
	if breach == nil {
		t.Fatal("expected breach after 4 flips")
	}
	if breach.Port != 443 {
		t.Errorf("expected port 443, got %d", breach.Port)
	}
	if breach.Flips != 4 {
		t.Errorf("expected 4 flips, got %d", breach.Flips)
	}
}

func TestRecord_OldEventsExpire(t *testing.T) {
	w := New(2, time.Minute)
	now := frozen()
	// record 2 flips far in the past
	w.Record(22, now.Add(-5*time.Minute))
	w.Record(22, now.Add(-4*time.Minute))
	// one fresh flip — should NOT breach because old ones are outside window
	b := w.Record(22, now)
	if b != nil {
		t.Fatalf("expected no breach after old events expired, got %+v", b)
	}
}

func TestReset_ClearsHistory(t *testing.T) {
	w := New(2, time.Minute)
	now := frozen()
	w.Record(8080, now)
	w.Record(8080, now.Add(time.Second))
	w.Reset(8080)
	b := w.Record(8080, now.Add(2*time.Second))
	if b != nil {
		t.Fatalf("expected no breach after reset, got %+v", b)
	}
}

func TestActive_ReturnsLiveFlips(t *testing.T) {
	w := New(10, time.Minute)
	now := frozen()
	w.Record(80, now.Add(-30*time.Second))
	w.Record(80, now.Add(-10*time.Second))
	w.Record(9000, now.Add(-90*time.Second)) // outside window
	active := w.Active(now)
	if active[80] != 2 {
		t.Errorf("expected 2 active flips for port 80, got %d", active[80])
	}
	if _, ok := active[9000]; ok {
		t.Error("port 9000 should not appear in active map")
	}
}

func TestNew_Defaults(t *testing.T) {
	w := New(0, 0)
	if w.threshold != 3 {
		t.Errorf("expected default threshold 3, got %d", w.threshold)
	}
	if w.window != 5*time.Minute {
		t.Errorf("expected default window 5m, got %v", w.window)
	}
}
