package rollup

import "github.com/user/portwatch/internal/state"

// Summary holds counts derived from a flushed batch of diffs.
type Summary struct {
	Opened []state.Diff
	Closed []state.Diff
}

// Summarise partitions a batch of diffs into opened and closed slices.
func Summarise(diffs []state.Diff) Summary {
	var s Summary
	for _, d := range diffs {
		switch d.Type {
		case state.DiffOpened:
			s.Opened = append(s.Opened, d)
		case state.DiffClosed:
			s.Closed = append(s.Closed, d)
		}
	}
	return s
}

// TotalChanged returns the total number of port changes in the summary.
func (s Summary) TotalChanged() int {
	return len(s.Opened) + len(s.Closed)
}

// HasChanges reports whether the summary contains any diffs.
func (s Summary) HasChanges() bool {
	return s.TotalChanged() > 0
}
