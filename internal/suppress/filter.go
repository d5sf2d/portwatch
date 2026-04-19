package suppress

import (
	"time"

	"github.com/user/portwatch/internal/state"
)

// FilterDiffs removes diffs whose port is currently suppressed.
// It returns only the diffs that should trigger alerts.
func FilterDiffs(diffs []state.Diff, l *List, t time.Time) []state.Diff {
	if l == nil {
		return diffs
	}
	out := diffs[:0:0]
	for _, d := range diffs {
		if !l.IsSuppressed(d.Port, t) {
			out = append(out, d)
		}
	}
	return out
}

// CountSuppressed returns the number of diffs that would be suppressed
// at the given time. This is useful for metrics and logging.
func CountSuppressed(diffs []state.Diff, l *List, t time.Time) int {
	if l == nil {
		return 0
	}
	count := 0
	for _, d := range diffs {
		if l.IsSuppressed(d.Port, t) {
			count++
		}
	}
	return count
}
