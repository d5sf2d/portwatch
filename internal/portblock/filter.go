package portblock

import "github.com/user/portwatch/internal/state"

// FilterDiffs removes any Diff entries whose port appears in the blocked
// registry. If reg is nil the original slice is returned unchanged.
func FilterDiffs(diffs []state.Diff, reg *Registry) []state.Diff {
	if reg == nil || len(diffs) == 0 {
		return diffs
	}
	out := diffs[:0:0]
	for _, d := range diffs {
		if !reg.IsBlocked(d.Port) {
			out = append(out, d)
		}
	}
	return out
}

// CountBlocked returns the number of diffs that would be filtered out by
// the given registry.
func CountBlocked(diffs []state.Diff, reg *Registry) int {
	if reg == nil {
		return 0
	}
	n := 0
	for _, d := range diffs {
		if reg.IsBlocked(d.Port) {
			n++
		}
	}
	return n
}
