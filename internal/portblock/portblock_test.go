package portblock_test

import (
	"testing"

	"github.com/user/portwatch/internal/portblock"
	"github.com/user/portwatch/internal/state"
)

func makeRegistry(t *testing.T, ports ...int) *portblock.Registry {
	t.Helper()
	reg := portblock.New()
	for _, p := range ports {
		if err := reg.Block(p, "test"); err != nil {
			t.Fatalf("Block(%d): %v", p, err)
		}
	}
	return reg
}

func makeDiff(port int, kind string) state.Diff {
	return state.Diff{Port: port, Kind: kind}
}

func TestBlock_InvalidPort(t *testing.T) {
	reg := portblock.New()
	if err := reg.Block(0, "bad"); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := reg.Block(65536, "bad"); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestIsBlocked_TrueAfterBlock(t *testing.T) {
	reg := makeRegistry(t, 22, 80)
	if !reg.IsBlocked(22) {
		t.Error("port 22 should be blocked")
	}
	if !reg.IsBlocked(80) {
		t.Error("port 80 should be blocked")
	}
	if reg.IsBlocked(443) {
		t.Error("port 443 should not be blocked")
	}
}

func TestUnblock_ClearsEntry(t *testing.T) {
	reg := makeRegistry(t, 22)
	reg.Unblock(22)
	if reg.IsBlocked(22) {
		t.Error("port 22 should no longer be blocked after Unblock")
	}
}

func TestReason_ReturnsReason(t *testing.T) {
	reg := portblock.New()
	_ = reg.Block(22, "ssh disabled")
	if got := reg.Reason(22); got != "ssh disabled" {
		t.Errorf("Reason(22) = %q; want %q", got, "ssh disabled")
	}
	if got := reg.Reason(80); got != "" {
		t.Errorf("Reason(80) = %q; want empty string", got)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	reg := makeRegistry(t, 22, 443)
	all := reg.All()
	if len(all) != 2 {
		t.Fatalf("All() len = %d; want 2", len(all))
	}
	// Mutating the copy must not affect the registry.
	delete(all, 22)
	if !reg.IsBlocked(22) {
		t.Error("mutating All() result should not affect registry")
	}
}

func TestFilterDiffs_BlockedPortRemoved(t *testing.T) {
	reg := makeRegistry(t, 22)
	diffs := []state.Diff{makeDiff(22, "opened"), makeDiff(80, "opened")}
	got := portblock.FilterDiffs(diffs, reg)
	if len(got) != 1 || got[0].Port != 80 {
		t.Errorf("FilterDiffs = %v; want only port 80", got)
	}
}

func TestFilterDiffs_NilRegistry(t *testing.T) {
	diffs := []state.Diff{makeDiff(22, "opened")}
	got := portblock.FilterDiffs(diffs, nil)
	if len(got) != 1 {
		t.Errorf("FilterDiffs with nil registry should return all diffs")
	}
}

func TestCountBlocked_Counts(t *testing.T) {
	reg := makeRegistry(t, 22, 23)
	diffs := []state.Diff{makeDiff(22, "opened"), makeDiff(23, "closed"), makeDiff(80, "opened")}
	if n := portblock.CountBlocked(diffs, reg); n != 2 {
		t.Errorf("CountBlocked = %d; want 2", n)
	}
}
