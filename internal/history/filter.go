package history

import "time"

// FilterOptions controls which entries are returned by Filter.
type FilterOptions struct {
	Since time.Time
	Until time.Time
	Port  int    // 0 means any
	Event string // "" means any
}

// Filter returns entries matching all non-zero criteria in opts.
func Filter(entries []Entry, opts FilterOptions) []Entry {
	var out []Entry
	for _, e := range entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if !opts.Until.IsZero() && e.Timestamp.After(opts.Until) {
			continue
		}
		if opts.Port != 0 && e.Port != opts.Port {
			continue
		}
		if opts.Event != "" && e.Event != opts.Event {
			continue
		}
		out = append(out, e)
	}
	return out
}
