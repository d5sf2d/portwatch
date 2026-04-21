package probe

import "time"

// Summary aggregates a slice of probe Results.
type Summary struct {
	Total      int
	Open       int
	Closed     int
	AvgLatency time.Duration
	MaxLatency time.Duration
}

// Summarise computes aggregate statistics over results.
func Summarise(results []Result) Summary {
	if len(results) == 0 {
		return Summary{}
	}
	var totalLatency time.Duration
	var maxLatency time.Duration
	open := 0
	for _, r := range results {
		if r.Open {
			open++
		}
		totalLatency += r.Latency
		if r.Latency > maxLatency {
			maxLatency = r.Latency
		}
	}
	return Summary{
		Total:      len(results),
		Open:       open,
		Closed:     len(results) - open,
		AvgLatency: totalLatency / time.Duration(len(results)),
		MaxLatency: maxLatency,
	}
}

// OpenResults returns only the results where the port was reachable.
func OpenResults(results []Result) []Result {
	out := make([]Result, 0, len(results))
	for _, r := range results {
		if r.Open {
			out = append(out, r)
		}
	}
	return out
}
