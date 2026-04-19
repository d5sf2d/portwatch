package suppress

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/state"
)

func makeDiff(port int, kind string) state.Diff {
	return state.Diff{Port: port, Kind: kind}
}

func TestFilterDiffs_SuppressedPortRemoved(t *testing.T) {
	l := New()
	t0 := time.Now()
	l.Add(Window{Port: 8080, Start: t0.Add(-time.Hour), End: t0.Add(time.Hour)})

	diffs := []state.Diff{makeDiff(8080, "opened"), makeDiff(443, "opened")}
	result := FilterDiffs(diffs, l, t0)

	if len(result) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(result))
	}
	if result[0].Port != 443 {
		t.Errorf("expected port 443, got %d", result[0].Port)
	}
}

func TestFilterDiffs_NilList(t *testing.T) {
	diffs := []state.Diff{makeDiff(8080, "opened")}
	result := FilterDiffs(diffs, nil, time.Now())
	if len(result) != 1 {
		t.Fatalf("expected 1 diff with nil list, got %d", len(result))
	}
}

func TestFilterDiffs_EmptyDiffs(t *testing.T) {
	l := New()
	result := FilterDiffs([]state.Diff{}, l, time.Now())
	if len(result) != 0 {
		t.Fatalf("expected 0 diffs, got %d", len(result))
	}
}
