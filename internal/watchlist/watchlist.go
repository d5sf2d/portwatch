// Package watchlist manages named port groups for organized monitoring.
package watchlist

import (
	"fmt"
	"sort"
)

// Group represents a named collection of ports.
type Group struct {
	Name  string `json:"name"`
	Ports []int  `json:"ports"`
}

// Watchlist holds multiple named port groups.
type Watchlist struct {
	Groups []Group `json:"groups"`
}

// AllPorts returns a deduplicated, sorted slice of all ports across all groups.
func (w *Watchlist) AllPorts() []int {
	seen := make(map[int]struct{})
	for _, g := range w.Groups {
		for _, p := range g.Ports {
			seen[p] = struct{}{}
		}
	}
	ports := make([]int, 0, len(seen))
	for p := range seen {
		ports = append(ports, p)
	}
	sort.Ints(ports)
	return ports
}

// Validate checks that group names are non-empty and ports are in valid range.
func (w *Watchlist) Validate() error {
	names := make(map[string]struct{})
	for _, g := range w.Groups {
		if g.Name == "" {
			return fmt.Errorf("watchlist: group name must not be empty")
		}
		if _, dup := names[g.Name]; dup {
			return fmt.Errorf("watchlist: duplicate group name %q", g.Name)
		}
		names[g.Name] = struct{}{}
		for _, p := range g.Ports {
			if p < 1 || p > 65535 {
				return fmt.Errorf("watchlist: port %d in group %q is out of range", p, g.Name)
			}
		}
	}
	return nil
}

// FindGroups returns the names of groups that contain the given port.
func (w *Watchlist) FindGroups(port int) []string {
	var names []string
	for _, g := range w.Groups {
		for _, p := range g.Ports {
			if p == port {
				names = append(names, g.Name)
				break
			}
		}
	}
	return names
}
