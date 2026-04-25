package porttrend

import (
	"testing"
	"time"
)

func frozenTracker(window time.Duration) *Tracker {
	return New(window)
}

func TestRecord_AppearsInTrends(t *testing.T) {
	tr := frozenTracker(time.Minute)
	tr.Record(Event{Port: 80, Host: "localhost", OpenedAt: time.Now()})

	summaries := tr.Trends(time.Now())
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}
	if summaries[0].Port != 80 {
		t.Errorf("expected port 80, got %d", summaries[0].Port)
	}
	if summaries[0].OpenCount != 1 {
		t.Errorf("expected OpenCount 1, got %d", summaries[0].OpenCount)
	}
}

func TestRecord_CountsOpenAndClose(t *testing.T) {
	tr := frozenTracker(time.Minute)
	now := time.Now()
	tr.Record(Event{Port: 443, Host: "host1", OpenedAt: now})
	tr.Record(Event{Port: 443, Host: "host1", OpenedAt: now.Add(time.Second), WasClosed: true})
	tr.Record(Event{Port: 443, Host: "host1", OpenedAt: now.Add(2 * time.Second)})

	summaries := tr.Trends(now.Add(3 * time.Second))
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}
	s := summaries[0]
	if s.OpenCount != 2 {
		t.Errorf("expected OpenCount 2, got %d", s.OpenCount)
	}
	if s.CloseCount != 1 {
		t.Errorf("expected CloseCount 1, got %d", s.CloseCount)
	}
}

func TestRecord_OldEventsExpire(t *testing.T) {
	tr := frozenTracker(5 * time.Second)
	old := time.Now().Add(-10 * time.Second)
	tr.Record(Event{Port: 22, Host: "h", OpenedAt: old})

	summaries := tr.Trends(time.Now())
	if len(summaries) != 0 {
		t.Errorf("expected 0 summaries after expiry, got %d", len(summaries))
	}
}

func TestRecord_MultipleHosts(t *testing.T) {
	tr := frozenTracker(time.Minute)
	now := time.Now()
	tr.Record(Event{Port: 80, Host: "host-a", OpenedAt: now})
	tr.Record(Event{Port: 80, Host: "host-b", OpenedAt: now})

	summaries := tr.Trends(now)
	if len(summaries) != 2 {
		t.Errorf("expected 2 summaries (one per host), got %d", len(summaries))
	}
}

func TestUnstable_ThresholdFilter(t *testing.T) {
	tr := frozenTracker(time.Minute)
	now := time.Now()
	for i := 0; i < 5; i++ {
		tr.Record(Event{Port: 8080, Host: "srv", OpenedAt: now.Add(time.Duration(i) * time.Second)})
	}
	tr.Record(Event{Port: 9090, Host: "srv", OpenedAt: now})

	results := Unstable(tr, 4)
	if len(results) != 1 {
		t.Fatalf("expected 1 unstable port, got %d", len(results))
	}
	if results[0].Port != 8080 {
		t.Errorf("expected port 8080 flagged, got %d", results[0].Port)
	}
}

func TestUnstable_NilTrackerReturnsNil(t *testing.T) {
	results := Unstable(nil, 1)
	if results != nil {
		t.Errorf("expected nil for nil tracker, got %v", results)
	}
}
