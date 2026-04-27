package portreport

import (
	"testing"
	"time"
)

func frozenTracker(ts time.Time) *Tracker {
	t := New()
	t.now = func() time.Time { return ts }
	return t
}

func TestRecordOpen_CreatesEntry(t *testing.T) {
	base := time.Now()
	tr := frozenTracker(base)
	tr.RecordOpen("localhost", 22, "ssh")
	r := tr.Report()
	if len(r.Stats) != 1 {
		t.Fatalf("expected 1 stat, got %d", len(r.Stats))
	}
	s := r.Stats[0]
	if s.Port != 22 || s.Host != "localhost" {
		t.Errorf("unexpected entry: %+v", s)
	}
	if s.OpenCount != 1 {
		t.Errorf("expected OpenCount 1, got %d", s.OpenCount)
	}
	if s.Label != "ssh" {
		t.Errorf("expected label ssh, got %q", s.Label)
	}
	if !s.FirstSeen.Equal(base) {
		t.Errorf("unexpected FirstSeen: %v", s.FirstSeen)
	}
}

func TestRecordOpen_AccumulatesCount(t *testing.T) {
	tr := New()
	tr.RecordOpen("host", 80, "http")
	tr.RecordOpen("host", 80, "http")
	tr.RecordOpen("host", 80, "http")
	r := tr.Report()
	if r.Stats[0].OpenCount != 3 {
		t.Errorf("expected 3 opens, got %d", r.Stats[0].OpenCount)
	}
}

func TestRecordClose_AccumulatesCount(t *testing.T) {
	tr := New()
	tr.RecordOpen("host", 443, "https")
	tr.RecordClose("host", 443)
	tr.RecordClose("host", 443)
	r := tr.Report()
	if r.Stats[0].CloseCount != 2 {
		t.Errorf("expected 2 closes, got %d", r.Stats[0].CloseCount)
	}
}

func TestReport_SortedByHostThenPort(t *testing.T) {
	tr := New()
	tr.RecordOpen("zhost", 22, "")
	tr.RecordOpen("ahost", 9000, "")
	tr.RecordOpen("ahost", 80, "")
	r := tr.Report()
	if r.Stats[0].Host != "ahost" || r.Stats[0].Port != 80 {
		t.Errorf("first entry should be ahost:80, got %s:%d", r.Stats[0].Host, r.Stats[0].Port)
	}
	if r.Stats[1].Port != 9000 {
		t.Errorf("second entry should be port 9000, got %d", r.Stats[1].Port)
	}
	if r.Stats[2].Host != "zhost" {
		t.Errorf("third entry should be zhost, got %s", r.Stats[2].Host)
	}
}

func TestReset_ClearsAllEntries(t *testing.T) {
	tr := New()
	tr.RecordOpen("host", 22, "ssh")
	tr.RecordOpen("host", 80, "http")
	tr.Reset()
	r := tr.Report()
	if len(r.Stats) != 0 {
		t.Errorf("expected empty report after reset, got %d entries", len(r.Stats))
	}
}

func TestReport_GeneratedAtIsSet(t *testing.T) {
	base := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	tr := frozenTracker(base)
	tr.RecordOpen("host", 8080, "")
	r := tr.Report()
	if !r.GeneratedAt.Equal(base) {
		t.Errorf("expected GeneratedAt %v, got %v", base, r.GeneratedAt)
	}
}
