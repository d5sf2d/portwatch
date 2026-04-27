package portage

import "github.com/user/portwatch/internal/state"

// FromDiffs updates the tracker based on a set of state diffs.
// Opened ports are observed; closed ports are forgotten.
func FromDiffs(t *Tracker, host string, diffs []state.Diff) {
	if t == nil {
		return
	}
	for _, d := range diffs {
		switch d.Type {
		case state.DiffOpened:
			t.Observe(host, d.Port)
		case state.DiffClosed:
			t.Forget(host, d.Port)
		}
	}
}

// Stale returns all records whose age exceeds the given threshold.
func Stale(t *Tracker, threshold int64) []Record {
	var out []Record
	for _, r := range t.All() {
		if int64(r.Age.Seconds()) >= threshold {
			out = append(out, r)
		}
	}
	return out
}
