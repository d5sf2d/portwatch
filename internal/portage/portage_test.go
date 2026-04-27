package portage

import (
	"testing"
	"time"
)

func frozenTracker(ts time.Time) *Tracker {
	t := New()
	t.now = func() time.Time { return ts }
	return t
}

func TestObserve_SetsFirstSeen(t *testing.T) {
	now := time.Now()
	tr := frozenTracker(now)
	tr.Observe("localhost", 80)
	age, ok := tr.Age("localhost", 80)
	if !ok {
		t.Fatal("expected port to be tracked")
	}
	if age != 0 {
		t.Fatalf("expected age 0, got %v", age)
	}
}

func TestObserve_PreservesFirstSeen(t *testing.T) {
	base := time.Now()
	tr := frozenTracker(base)
	tr.Observe("localhost", 443)
	tr.now = func() time.Time { return base.Add(5 * time.Minute) }
	tr.Observe("localhost", 443) // second call should not reset
	age, ok := tr.Age("localhost", 443)
	if !ok {
		t.Fatal("expected port to be tracked")
	}
	if age != 5*time.Minute {
		t.Fatalf("expected 5m, got %v", age)
	}
}

func TestForget_RemovesPort(t *testing.T) {
	tr := frozenTracker(time.Now())
	tr.Observe("localhost", 22)
	tr.Forget("localhost", 22)
	_, ok := tr.Age("localhost", 22)
	if ok {
		t.Fatal("expected port to be forgotten")
	}
}

func TestAge_UnknownPort(t *testing.T) {
	tr := New()
	_, ok := tr.Age("localhost", 9999)
	if ok {
		t.Fatal("expected false for unknown port")
	}
}

func TestAll_ReturnsAllRecords(t *testing.T) {
	base := time.Now()
	tr := frozenTracker(base)
	tr.Observe("host1", 80)
	tr.Observe("host1", 443)
	tr.Observe("host2", 22)
	records := tr.All()
	if len(records) != 3 {
		t.Fatalf("expected 3 records, got %d", len(records))
	}
}

func TestStale_FiltersCorrectly(t *testing.T) {
	base := time.Now()
	tr := frozenTracker(base)
	tr.Observe("localhost", 80)
	tr.now = func() time.Time { return base.Add(2 * time.Hour) }
	tr.Observe("localhost", 443)
	// 443 was just added at +2h so age=0; 80 has age=2h
	stale := Stale(tr, int64((1 * time.Hour).Seconds()))
	if len(stale) != 1 {
		t.Fatalf("expected 1 stale record, got %d", len(stale))
	}
	if stale[0].Port != 80 {
		t.Fatalf("expected port 80, got %d", stale[0].Port)
	}
}
