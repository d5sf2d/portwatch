package debounce

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/state"
)

func makeDiff(host string, port int, change string) state.Diff {
	return state.Diff{Host: host, Port: port, Change: change}
}

func TestIsDuplicate_FirstCallIsNotDuplicate(t *testing.T) {
	db := New(5 * time.Second)
	d := makeDiff("localhost", 80, "opened")
	if db.IsDuplicate(d) {
		t.Fatal("first occurrence should not be a duplicate")
	}
}

func TestIsDuplicate_SecondCallWithinWindowIsDuplicate(t *testing.T) {
	db := New(5 * time.Second)
	d := makeDiff("localhost", 80, "opened")
	db.IsDuplicate(d)
	if !db.IsDuplicate(d) {
		t.Fatal("second call within window should be a duplicate")
	}
}

func TestIsDuplicate_AfterWindowExpiry(t *testing.T) {
	now := time.Now()
	db := New(5 * time.Second)
	db.nowFunc = func() time.Time { return now }

	d := makeDiff("localhost", 443, "closed")
	db.IsDuplicate(d)

	// advance past window
	db.nowFunc = func() time.Time { return now.Add(6 * time.Second) }
	if db.IsDuplicate(d) {
		t.Fatal("should not be a duplicate after window expires")
	}
}

func TestIsDuplicate_DifferentPortNotDuplicate(t *testing.T) {
	db := New(5 * time.Second)
	db.IsDuplicate(makeDiff("localhost", 80, "opened"))
	if db.IsDuplicate(makeDiff("localhost", 443, "opened")) {
		t.Fatal("different port should not be considered a duplicate")
	}
}

func TestFilter_RemovesDuplicates(t *testing.T) {
	db := New(5 * time.Second)
	diffs := []state.Diff{
		makeDiff("localhost", 80, "opened"),
		makeDiff("localhost", 443, "opened"),
	}

	// first pass — all pass through
	out := db.Filter(diffs)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}

	// second pass — all suppressed
	out = db.Filter(diffs)
	if len(out) != 0 {
		t.Fatalf("expected 0, got %d", len(out))
	}
}

func TestExpire_RemovesStaleEntries(t *testing.T) {
	now := time.Now()
	db := New(5 * time.Second)
	db.nowFunc = func() time.Time { return now }

	db.IsDuplicate(makeDiff("localhost", 22, "opened"))

	db.nowFunc = func() time.Time { return now.Add(10 * time.Second) }
	db.Expire()

	if len(db.seen) != 0 {
		t.Fatalf("expected empty seen map after expire, got %d entries", len(db.seen))
	}
}
