package portping

import "time"

// Summary aggregates a slice of ping Results.
type Summary struct {
	Total   int
	Alive   int
	Dead    int
	AvgRTT  time.Duration
	MaxRTT  time.Duration
	MinRTT  time.Duration
}

// Summarise computes aggregate statistics from a set of Results.
func Summarise(results []Result) Summary {
	if len(results) == 0 {
		return Summary{}
	}

	var s Summary
	s.Total = len(results)
	s.MinRTT = time.Duration(1<<63 - 1)

	var totalRTT time.Duration
	for _, r := range results {
		if r.Alive {
			s.Alive++
			totalRTT += r.Latency
			if r.Latency > s.MaxRTT {
				s.MaxRTT = r.Latency
			}
			if r.Latency < s.MinRTT {
				s.MinRTT = r.Latency
			}
		} else {
			s.Dead++
		}
	}

	if s.Alive > 0 {
		s.AvgRTT = totalRTT / time.Duration(s.Alive)
	} else {
		s.MinRTT = 0
	}
	return s
}

// AliveResults filters and returns only the results where Alive is true.
func AliveResults(results []Result) []Result {
	out := make([]Result, 0, len(results))
	for _, r := range results {
		if r.Alive {
			out = append(out, r)
		}
	}
	return out
}
