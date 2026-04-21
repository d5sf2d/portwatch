package watchdog

import (
	"time"

	"github.com/user/portwatch/internal/state"
)

// CheckDiffs feeds every diff entry into the watchdog and returns all
// breaches found.  Only Open and Closed transitions are considered flips;
// unchanged entries are ignored.
func CheckDiffs(w *Watchdog, diffs []state.Diff, at time.Time) []Breach {
	if w == nil || len(diffs) == 0 {
		return nil
	}
	var breaches []Breach
	for _, d := range diffs {
		if d.Status != state.StatusOpened && d.Status != state.StatusClosed {
			continue
		}
		if b := w.Record(d.Port, at); b != nil {
			breaches = append(breaches, *b)
		}
	}
	return breaches
}
