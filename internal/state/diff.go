package state

// Change describes a port that opened or closed between two snapshots.
type Change struct {
	Port     int
	Protocol string
	// Opened is true when the port transitioned closed→open.
	Opened bool
	// Closed is true when the port transitioned open→closed.
	Closed bool
}

// Diff compares a previous snapshot against a current one and returns
// any ports whose open/closed status changed.
func Diff(prev, curr Snapshot) []Change {
	prevMap := make(map[int]PortState, len(prev.Ports))
	for _, p := range prev.Ports {
		prevMap[p.Port] = p
	}

	currMap := make(map[int]PortState, len(curr.Ports))
	for _, p := range curr.Ports {
		currMap[p.Port] = p
	}

	var changes []Change

	// Detect newly opened or closed ports present in current scan.
	for _, cp := range curr.Ports {
		pp, existed := prevMap[cp.Port]
		if !existed && cp.Open {
			changes = append(changes, Change{Port: cp.Port, Protocol: cp.Protocol, Opened: true})
		} else if existed && pp.Open != cp.Open {
	(changes, Change{Port: cp.Port, Protocol: cp.Protocol, Opened: cp.Open, Closed: !cp.Open})
		}
	}

	// Detect ports that disappeared entirely ( as closed).
	for _, pp := range prev.Ports {
		if pp.Open {
			if _, found := currMap[pp.Port]; !found {
				changes = append(changes, Change{Port: pp.Port, Protocol: pp.Protocol, Closed: true})
			}
		}
	}

	return changes
}
