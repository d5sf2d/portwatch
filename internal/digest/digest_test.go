package digest_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/digest"
	"github.com/user/portwatch/internal/state"
)

func makeSnap(ports ...int) state.Snapshot {
	return state.Snapshot{Host: "localhost", Ports: ports, ScannedAt: time.Now()}
}

func TestOf_SamePortsSameHash(t *testing.T) {
	a := makeSnap(22, 80, 443)
	b := makeSnap(22, 80, 443)
	if digest.Of(a) != digest.Of(b) {
		t.Fatal("expected identical snapshots to produce the same digest")
	}
}

func TestOf_OrderIndependent(t *testing.T) {
	a := makeSnap(443, 80, 22)
	b := makeSnap(22, 80, 443)
	if digest.Of(a) != digest.Of(b) {
		t.Fatal("expected port order not to affect digest")
	}
}

func TestOf_DifferentPortsDifferentHash(t *testing.T) {
	a := makeSnap(22, 80)
	b := makeSnap(22, 443)
	if digest.Of(a) == digest.Of(b) {
		t.Fatal("expected different ports to produce different digests")
	}
}

func TestChanged_DetectsChange(t *testing.T) {
	prev := makeSnap(22, 80)
	curr := makeSnap(22, 80, 8080)
	if !digest.Changed(prev, curr) {
		t.Fatal("expected Changed to return true when a port was added")
	}
}

func TestChanged_NoChange(t *testing.T) {
	snap := makeSnap(22, 80)
	if digest.Changed(snap, snap) {
		t.Fatal("expected Changed to return false for identical snapshots")
	}
}

func TestCompare_Summary(t *testing.T) {
	prev := makeSnap(22)
	curr := makeSnap(22, 80)
	s := digest.Compare(prev, curr)
	if !s.Changed {
		t.Fatal("summary should report Changed=true")
	}
	if s.Prev == s.Curr {
		t.Fatal("prev and curr digests should differ")
	}
}

func TestDigest_Equal(t *testing.T) {
	snap := makeSnap(22, 80)
	d := digest.Of(snap)
	if !d.Equal(d) {
		t.Fatal("a digest must equal itself")
	}
}

func TestDigest_String_NonEmpty(t *testing.T) {
	d := digest.Of(makeSnap(80))
	if d.String() == "" {
		t.Fatal("digest string must not be empty")
	}
}
