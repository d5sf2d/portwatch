package history

import (
	"time"

	"github.com/example/portwatch/internal/state"
)

// FromDiffs converts a slice of state.Diff values into history Entries using
// the provided timestamp and host string.
func FromDiffs(diffs []state.Diff, host string, ts time.Time) []Entry {
	entries := make([]Entry, 0, len(diffs))
	for _, d := range diffs {
		event := "opened"
		if d.Type == state.Closed {
			event = "closed"
		}
		entries = append(entries, Entry{
			Timestamp: ts,
			Port:      d.Port,
			Proto:     d.Proto,
			Event:     event,
			Host:      host,
		})
	}
	return entries
}

// FilterByHost returns only the entries that match the given host string.
func FilterByHost(entries []Entry, host string) []Entry {
	result := make([]Entry, 0)
	for _, e := range entries {
		if e.Host == host {
			result = append(result, e)
		}
	}
	return result
}
