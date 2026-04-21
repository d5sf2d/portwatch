package digest

import "github.com/user/portwatch/internal/state"

// Summary describes the outcome of comparing two snapshots via their digests.
type Summary struct {
	Prev    Digest
	Curr    Digest
	Changed bool
}

// Compare builds a Summary for the transition from prev to curr.
func Compare(prev, curr state.Snapshot) Summary {
	p := Of(prev)
	c := Of(curr)
	return Summary{
		Prev:    p,
		Curr:    c,
		Changed: !p.Equal(c),
	}
}
