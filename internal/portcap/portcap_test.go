package portcap_test

import (
	"testing"

	"portwatch/internal/portcap"
)

func makeTracker(t *testing.T, defaultCap int) *portcap.Tracker {
	t.Helper()
	tr, err := portcap.New(defaultCap)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return tr
}

func TestNew_InvalidCap(t *testing.T) {
	_, err := portcap.New(0)
	if err == nil {
		t.Fatal("expected error for zero cap")
	}
	_, err = portcap.New(-5)
	if err == nil {
		t.Fatal("expected error for negative cap")
	}
}

func TestCap_DefaultApplied(t *testing.T) {
	tr := makeTracker(t, 10)
	if got := tr.Cap("host-a"); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestSetHost_OverridesDefault(t *testing.T) {
	tr := makeTracker(t, 10)
	if err := tr.SetHost("host-b", 3); err != nil {
		t.Fatalf("SetHost: %v", err)
	}
	if got := tr.Cap("host-b"); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
	// default unchanged for other hosts
	if got := tr.Cap("host-a"); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestSetHost_InvalidCap(t *testing.T) {
	tr := makeTracker(t, 10)
	if err := tr.SetHost("host-c", 0); err == nil {
		t.Fatal("expected error for zero host cap")
	}
}

func TestExceeded_WithinLimit(t *testing.T) {
	tr := makeTracker(t, 5)
	if tr.Exceeded("h", 5) {
		t.Fatal("expected not exceeded at exactly the cap")
	}
}

func TestExceeded_OverLimit(t *testing.T) {
	tr := makeTracker(t, 5)
	if !tr.Exceeded("h", 6) {
		t.Fatal("expected exceeded when openCount > cap")
	}
}

func TestOverage_BelowCap(t *testing.T) {
	tr := makeTracker(t, 10)
	if got := tr.Overage("h", 7); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestOverage_AboveCap(t *testing.T) {
	tr := makeTracker(t, 10)
	if got := tr.Overage("h", 13); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}
