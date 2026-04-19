package baseline

// Delta describes ports that deviate from the baseline.
type Delta struct {
	Added   []int // ports open now but not in baseline
	Removed []int // ports in baseline but no longer open
}

// HasChanges returns true if the delta is non-empty.
func (d Delta) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0
}

// Compare returns a Delta between the saved baseline and the current ports.
func Compare(b Baseline, current []int) Delta {
	baseSet := toSet(b.Ports)
	curSet := toSet(current)

	var added, removed []int

	for p := range curSet {
		if !baseSet[p] {
			added = append(added, p)
		}
	}
	for p := range baseSet {
		if !curSet[p] {
			removed = append(removed, p)
		}
	}
	return Delta{Added: added, Removed: removed}
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
