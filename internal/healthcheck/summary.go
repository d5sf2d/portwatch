package healthcheck

// Summary aggregates a slice of Status values into counts.
type Summary struct {
	Total     int
	Healthy   int
	Unhealthy int
}

// Summarise builds a Summary from a slice of Status results.
func Summarise(statuses []Status) Summary {
	s := Summary{Total: len(statuses)}
	for _, st := range statuses {
		if st.Healthy {
			s.Healthy++
		} else {
			s.Unhealthy++
		}
	}
	return s
}

// Unhealthy returns only the failed statuses from a slice.
func Unhealthy(statuses []Status) []Status {
	var out []Status
	for _, s := range statuses {
		if !s.Healthy {
			out = append(out, s)
		}
	}
	return out
}
