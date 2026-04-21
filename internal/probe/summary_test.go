package probe_test

import (
	"testing"
	"time"

	"portwatch/internal/probe"
)

func makeResults() []probe.Result {
	return []probe.Result{
		{Host: "h", Port: 80, Open: true, Latency: 10 * time.Millisecond},
		{Host: "h", Port: 443, Open: true, Latency: 20 * time.Millisecond},
		{Host: "h", Port: 8080, Open: false, Latency: 5 * time.Millisecond},
	}
}

func TestSummarise_Counts(t *testing.T) {
	s := probe.Summarise(makeResults())
	if s.Total != 3 {
		t.Errorf("total: want 3, got %d", s.Total)
	}
	if s.Open != 2 {
		t.Errorf("open: want 2, got %d", s.Open)
	}
	if s.Closed != 1 {
		t.Errorf("closed: want 1, got %d", s.Closed)
	}
}

func TestSummarise_Latency(t *testing.T) {
	s := probe.Summarise(makeResults())
	if s.MaxLatency != 20*time.Millisecond {
		t.Errorf("max latency: want 20ms, got %v", s.MaxLatency)
	}
	wantAvg := (10 + 20 + 5) * time.Millisecond / 3
	if s.AvgLatency != wantAvg {
		t.Errorf("avg latency: want %v, got %v", wantAvg, s.AvgLatency)
	}
}

func TestSummarise_Empty(t *testing.T) {
	s := probe.Summarise(nil)
	if s.Total != 0 || s.Open != 0 {
		t.Errorf("expected zero summary for nil input")
	}
}

func TestOpenResults_FiltersCorrectly(t *testing.T) {
	open := probe.OpenResults(makeResults())
	if len(open) != 2 {
		t.Fatalf("want 2 open results, got %d", len(open))
	}
	for _, r := range open {
		if !r.Open {
			t.Errorf("closed result leaked into OpenResults")
		}
	}
}

func TestOpenResults_AllClosed(t *testing.T) {
	results := []probe.Result{
		{Port: 1, Open: false},
		{Port: 2, Open: false},
	}
	open := probe.OpenResults(results)
	if len(open) != 0 {
		t.Errorf("expected empty slice, got %d", len(open))
	}
}
