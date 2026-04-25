package portmute

import "github.com/user/portwatch/internal/state"

// FilterDiffs removes diff entries whose port is currently muted.
// If list is nil, diffs are returned unchanged.
func FilterDiffs(diffs []state.Diff, list *List) []state.Diff {
	if list == nil || len(diffs) == 0 {
		return diffs
	}
	out := diffs[:0:0]
	for _, d := range diffs {
		if !list.IsMuted(d.Port) {
			out = append(out, d)
		}
	}
	return out
}

// CountMuted returns the number of diffs that would be suppressed by the list.
func CountMuted(diffs []state.Diff, list *List) int {
	if list == nil {
		return 0
	}
	count := 0
	for _, d := range diffs {
		if list.IsMuted(d.Port) {
			count++
		}
	}
	return count
}
