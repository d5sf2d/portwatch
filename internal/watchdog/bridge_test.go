package watchdog

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/state"
)

func makeDiff(port int, status state.DiffStatus) state.Diff {
	return state.Diff{Port: port, Status: status}
}

func TestCheckDiffs_DetectsFlap(t *testing.T) {
	w := New(2, time.Minute)
	now := time.Now()
	diffs := []state.Diff{
		makeDiff(80, state.StatusOpened),
		makeDiff(80, state.StatusClosed),
		makeDiff(80, state.StatusOpened),
	}
	var breaches []Breach
	for _, d := range diffs {
		breaches = CheckDiffs(w, []state.Diff{d}, now)
	}
	if len(breaches) == 0 {
		t.Fatal("expected at least one breach")
	}
	if breaches[0].Port != 80 {
		t.Errorf("expected breach on port 80, got %d", breaches[0].Port)
	}
}

func TestCheckDiffs_NilWatchdog(t *testing.T) {
	diffs := []state.Diff{makeDiff(443, state.StatusOpened)}
	breaches := CheckDiffs(nil, diffs, time.Now())
	if breaches != nil {
		t.Errorf("expected nil breaches for nil watchdog")
	}
}

func TestCheckDiffs_EmptyDiffs(t *testing.T) {
	w := New(2, time.Minute)
	breaches := CheckDiffs(w, nil, time.Now())
	if len(breaches) != 0 {
		t.Errorf("expected no breaches for empty diffs")
	}
}
