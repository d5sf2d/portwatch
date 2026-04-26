package portquota

import "github.com/example/portwatch/internal/state"

// ExceedingPorts returns the ports from snap that cause the open-port count to
// exceed the quota for snap.Host. The returned slice contains only the ports
// beyond the limit, ordered by port number (snap already stores ports sorted).
// If the quota is not exceeded, nil is returned.
func ExceedingPorts(q *Quota, snap state.Snapshot) []int {
	if q == nil {
		return nil
	}
	limit := q.Limit(snap.Host)
	if limit == 0 || len(snap.Ports) <= limit {
		return nil
	}
	// Ports beyond the limit are considered the "excess" set.
	excess := make([]int, 0, len(snap.Ports)-limit)
	for i, p := range snap.Ports {
		if i >= limit {
			excess = append(excess, p.Port)
		}
	}
	return excess
}

// CountExceeding returns the number of ports beyond the quota, or 0 if the
// quota is not exceeded.
func CountExceeding(q *Quota, snap state.Snapshot) int {
	excess := ExceedingPorts(q, snap)
	return len(excess)
}
