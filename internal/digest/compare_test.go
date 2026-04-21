package digest_test

import (
	"testing"

	"github.com/user/portwatch/internal/digest"
)

func TestCompare_Unchanged(t *testing.T) {
	snap := makeSnap(22, 80, 443)
	s := digest.Compare(snap, snap)
	if s.Changed {
		t.Fatal("expected Changed=false for identical snapshots")
	}
	if s.Prev != s.Curr {
		t.Fatal("expected Prev and Curr digests to be equal")
	}
}

func TestCompare_EmptySnapshots(t *testing.T) {
	a := makeSnap()
	b := makeSnap()
	s := digest.Compare(a, b)
	if s.Changed {
		t.Fatal("two empty snapshots should not differ")
	}
}

func TestCompare_FromEmptyToNonEmpty(t *testing.T) {
	empty := makeSnap()
	full := makeSnap(22, 80)
	s := digest.Compare(empty, full)
	if !s.Changed {
		t.Fatal("expected Changed=true when ports appear")
	}
}

func TestCompare_HashLengthIs64(t *testing.T) {
	s := digest.Compare(makeSnap(80), makeSnap(443))
	if len(s.Curr.Hash) != 64 {
		t.Fatalf("expected 64-char hex hash, got %d", len(s.Curr.Hash))
	}
}

func TestCompare_PrevAndCurrDifferWhenChanged(t *testing.T) {
	prev := makeSnap(22, 80)
	curr := makeSnap(22, 80, 443)
	s := digest.Compare(prev, curr)
	if !s.Changed {
		t.Fatal("expected Changed=true when port set differs")
	}
	if s.Prev == s.Curr {
		t.Fatal("expected Prev and Curr digests to differ when snapshots differ")
	}
}
