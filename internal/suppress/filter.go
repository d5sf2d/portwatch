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
