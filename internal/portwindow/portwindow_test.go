package portwindow

import (
	"testing"
	"time"
)

func frozenWindow(d time.Duration, now time.Time) *Window {
	w := New(d)
	w.now = func() time.Time { return now }
	return w
}

func TestWindow_Opens(t *testing.T) {
	now := time.Now()
	w := frozenWindow(10*time.Second, now)
	w.Record(true)
	w.Record(true)
	w.Record(false)
	if got := w.Opens(); got != 2 {
		t.Fatalf("Opens() = %d, want 2", got)
	}
}

func TestWindow_Closes(t *testing.T) {
	now := time.Now()
	w := frozenWindow(10*time.Second, now)
	w.Record(false)
	w.Record(true)
	if got := w.Closes(); got != 1 {
		t.Fatalf("Closes() = %d, want 1", got)
	}
}

func TestWindow_Total(t *testing.T) {
	now := time.Now()
	w := frozenWindow(10*time.Second, now)
	w.Record(true)
	w.Record(false)
	w.Record(true)
	if got := w.Total(); got != 3 {
		t.Fatalf("Total() = %d, want 3", got)
	}
}

func TestWindow_EventsExpire(t *testing.T) {
	base := time.Now()
	current := base
	w := New(5 * time.Second)
	w.now = func() time.Time { return current }

	w.Record(true) // at base
	current = base.Add(6 * time.Second)
	w.Record(false) // at base+6s

	// Only the second event should remain.
	if got := w.Total(); got != 1 {
		t.Fatalf("Total() after expiry = %d, want 1", got)
	}
	if got := w.Closes(); got != 1 {
		t.Fatalf("Closes() after expiry = %d, want 1", got)
	}
}

func TestRegistry_RecordAndGet(t *testing.T) {
	r := NewRegistry(30 * time.Second)
	r.Record("localhost", 22, true)
	r.Record("localhost", 22, false)

	w := r.Get("localhost", 22)
	if w == nil {
		t.Fatal("expected window, got nil")
	}
	if got := w.Total(); got != 2 {
		t.Fatalf("Total() = %d, want 2", got)
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	r := NewRegistry(30 * time.Second)
	if got := r.Get("localhost", 9999); got != nil {
		t.Fatalf("expected nil for unseen key, got %v", got)
	}
}

func TestRegistry_ActiveKeys(t *testing.T) {
	r := NewRegistry(30 * time.Second)
	r.Record("host1", 80, true)
	r.Record("host2", 443, false)

	keys := r.ActiveKeys()
	if len(keys) != 2 {
		t.Fatalf("ActiveKeys() len = %d, want 2", len(keys))
	}
}
