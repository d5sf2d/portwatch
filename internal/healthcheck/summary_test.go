package healthcheck_test

import (
	"testing"

	"github.com/user/portwatch/internal/healthcheck"
)

func makeStatuses() []healthcheck.Status {
	return []healthcheck.Status{
		{Host: "localhost", Port: 80, Healthy: true, Reason: "ok"},
		{Host: "localhost", Port: 443, Healthy: true, Reason: "ok"},
		{Host: "localhost", Port: 9999, Healthy: false, Reason: "refused"},
	}
}

func TestSummarise_Counts(t *testing.T) {
	s := healthcheck.Summarise(makeStatuses())

	if s.Total != 3 {
		t.Errorf("total: got %d want 3", s.Total)
	}
	if s.Healthy != 2 {
		t.Errorf("healthy: got %d want 2", s.Healthy)
	}
	if s.Unhealthy != 1 {
		t.Errorf("unhealthy: got %d want 1", s.Unhealthy)
	}
}

func TestSummarise_Empty(t *testing.T) {
	s := healthcheck.Summarise(nil)
	if s.Total != 0 || s.Healthy != 0 || s.Unhealthy != 0 {
		t.Errorf("expected all zeros, got %+v", s)
	}
}

func TestUnhealthy_FiltersCorrectly(t *testing.T) {
	bad := healthcheck.Unhealthy(makeStatuses())
	if len(bad) != 1 {
		t.Fatalf("expected 1 unhealthy, got %d", len(bad))
	}
	if bad[0].Port != 9999 {
		t.Errorf("expected port 9999, got %d", bad[0].Port)
	}
}

func TestUnhealthy_AllHealthy(t *testing.T) {
	statuses := []healthcheck.Status{
		{Healthy: true},
		{Healthy: true},
	}
	if got := healthcheck.Unhealthy(statuses); len(got) != 0 {
		t.Errorf("expected empty slice, got %d items", len(got))
	}
}
