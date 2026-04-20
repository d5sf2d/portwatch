package metrics_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestRecordScan_IncrementsCounter(t *testing.T) {
	c := metrics.New()
	c.RecordScan(50*time.Millisecond, 3)

	snap := c.Snapshot()
	if snap.ScansTotal != 1 {
		t.Fatalf("expected ScansTotal=1, got %d", snap.ScansTotal)
	}
	if snap.PortsOpen != 3 {
		t.Fatalf("expected PortsOpen=3, got %d", snap.PortsOpen)
	}
	if snap.LastScanDur != 50*time.Millisecond {
		t.Fatalf("unexpected LastScanDur: %v", snap.LastScanDur)
	}
	if snap.LastScanAt.IsZero() {
		t.Fatal("LastScanAt should not be zero after RecordScan")
	}
}

func TestRecordAlert_Accumulates(t *testing.T) {
	c := metrics.New()
	c.RecordAlert(2)
	c.RecordAlert(3)

	snap := c.Snapshot()
	if snap.AlertsTotal != 5 {
		t.Fatalf("expected AlertsTotal=5, got %d", snap.AlertsTotal)
	}
}

func TestRecordSuppressed_Accumulates(t *testing.T) {
	c := metrics.New()
	c.RecordSuppressed(1)
	c.RecordSuppressed(4)

	snap := c.Snapshot()
	if snap.SuppressedTotal != 5 {
		t.Fatalf("expected SuppressedTotal=5, got %d", snap.SuppressedTotal)
	}
}

func TestReset_ZeroesCounters(t *testing.T) {
	c := metrics.New()
	c.RecordScan(10*time.Millisecond, 7)
	c.RecordAlert(3)
	c.RecordSuppressed(2)
	c.Reset()

	snap := c.Snapshot()
	if snap.ScansTotal != 0 || snap.AlertsTotal != 0 || snap.SuppressedTotal != 0 || snap.PortsOpen != 0 {
		t.Fatalf("expected all counters zero after Reset, got %+v", snap)
	}
	if !snap.LastScanAt.IsZero() {
		t.Fatal("LastScanAt should be zero after Reset")
	}
}

func TestSnapshot_IsConsistentCopy(t *testing.T) {
	c := metrics.New()
	c.RecordScan(5*time.Millisecond, 2)

	s1 := c.Snapshot()
	c.RecordScan(5*time.Millisecond, 9)
	s2 := c.Snapshot()

	if s1.ScansTotal != 1 {
		t.Fatalf("s1 should be independent of later mutations, got ScansTotal=%d", s1.ScansTotal)
	}
	if s2.ScansTotal != 2 {
		t.Fatalf("s2 ScansTotal expected 2, got %d", s2.ScansTotal)
	}
	if s2.PortsOpen != 9 {
		t.Fatalf("s2 PortsOpen expected 9, got %d", s2.PortsOpen)
	}
}
