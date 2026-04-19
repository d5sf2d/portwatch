package baseline_test

import (
	"testing"

	"portwatch/internal/baseline"
)

func TestDelta_HasChanges_True(t *testing.T) {
	d := baseline.Delta{Added: []int{8080}}
	if !d.HasChanges() {
		t.Error("expected HasChanges true")
	}
}

func TestDelta_HasChanges_False(t *testing.T) {
	d := baseline.Delta{}
	if d.HasChanges() {
		t.Error("expected HasChanges false")
	}
}

func TestCompare_BothDirections(t *testing.T) {
	b := baseline.Capture([]int{22, 80, 443})
	current := []int{80, 443, 3306}
	d := baseline.Compare(b, current)

	if len(d.Added) != 1 {
		t.Errorf("expected 1 added, got %d", len(d.Added))
	}
	if len(d.Removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(d.Removed))
	}
}

func TestCompare_EmptyBaseline(t *testing.T) {
	b := baseline.Capture([]int{})
	d := baseline.Compare(b, []int{80})
	if len(d.Added) != 1 {
		t.Errorf("expected 1 added, got %d", len(d.Added))
	}
}

func TestCompare_EmptyCurrent(t *testing.T) {
	b := baseline.Capture([]int{80, 443})
	d := baseline.Compare(b, []int{})
	if len(d.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(d.Removed))
	}
}
