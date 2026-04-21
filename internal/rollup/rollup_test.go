package rollup

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/state"
)

func makeDiff(port int, typ state.DiffType) state.Diff {
	return state.Diff{Port: port, Type: typ}
}

func TestAdd_BuffersWithinWindow(t *testing.T) {
	r := New(10*time.Second, 5)
	batch, ready := r.Add(makeDiff(80, state.DiffOpened))
	if ready {
		t.Fatal("expected not ready after first add")
	}
	if batch != nil {
		t.Fatal("expected nil batch")
	}
	if r.Pending() != 1 {
		t.Fatalf("expected 1 pending, got %d", r.Pending())
	}
}

func TestAdd_FlushesAtCapacity(t *testing.T) {
	r := New(10*time.Second, 3)
	for i := 0; i < 2; i++ {
		_, ready := r.Add(makeDiff(80+i, state.DiffOpened))
		if ready {
			t.Fatal("flushed too early")
		}
	}
	batch, ready := r.Add(makeDiff(443, state.DiffClosed))
	if !ready {
		t.Fatal("expected flush at capacity")
	}
	if len(batch) != 3 {
		t.Fatalf("expected 3 diffs, got %d", len(batch))
	}
	if r.Pending() != 0 {
		t.Fatal("buffer should be empty after flush")
	}
}

func TestAdd_FlushesAfterWindowExpiry(t *testing.T) {
	r := New(100*time.Millisecond, 10)
	past := time.Now().Add(-200 * time.Millisecond)
	r.now = func() time.Time { return past }
	r.Add(makeDiff(22, state.DiffOpened)) // sets deadline in the past
	r.now = time.Now

	batch, ready := r.Add(makeDiff(8080, state.DiffOpened))
	if !ready {
		t.Fatal("expected flush after window expiry")
	}
	if len(batch) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(batch))
	}
}

func TestFlush_EmptyReturnsNil(t *testing.T) {
	r := New(time.Second, 5)
	if r.Flush() != nil {
		t.Fatal("expected nil flush on empty buffer")
	}
}

func TestSummarise_PartitionsDiffs(t *testing.T) {
	diffs := []state.Diff{
		makeDiff(80, state.DiffOpened),
		makeDiff(443, state.DiffOpened),
		makeDiff(22, state.DiffClosed),
	}
	s := Summarise(diffs)
	if len(s.Opened) != 2 {
		t.Fatalf("expected 2 opened, got %d", len(s.Opened))
	}
	if len(s.Closed) != 1 {
		t.Fatalf("expected 1 closed, got %d", len(s.Closed))
	}
	if !s.HasChanges() {
		t.Fatal("expected HasChanges true")
	}
	if s.TotalChanged() != 3 {
		t.Fatalf("expected total 3, got %d", s.TotalChanged())
	}
}

func TestSummarise_Empty(t *testing.T) {
	s := Summarise(nil)
	if s.HasChanges() {
		t.Fatal("expected no changes for empty input")
	}
}
