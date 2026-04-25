package porttrend

import (
	"time"

	"github.com/user/portwatch/internal/state"
)

// FromDiffs converts a slice of state.Diff values into Events and records
// them in the Tracker. host identifies the scanned target.
func FromDiffs(t *Tracker, host string, diffs []state.Diff) {
	if t == nil {
		return
	}
	now := time.Now()
	for _, d := range diffs {
		e := Event{
			Port:      d.Port,
			Host:      host,
			OpenedAt:  now,
			WasClosed: d.Status == state.Closed,
		}
		t.Record(e)
	}
}

// Unstable returns summaries for ports whose total event count (opens +
// closes) meets or exceeds the threshold within the tracker's window.
func Unstable(t *Tracker, threshold int) []Summary {
	if t == nil {
		return nil
	}
	all := t.Trends(time.Now())
	out := make([]Summary, 0)
	for _, s := range all {
		if s.OpenCount+s.CloseCount >= threshold {
			out = append(out, s)
		}
	}
	return out
}
