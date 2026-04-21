// Package grouping organises ports into named logical groups and provides
// helpers to query membership and retrieve group metadata.
package grouping

import (
	"errors"
	"fmt"
)

// Group is a named collection of ports.
type Group struct {
	Name  string
	Ports []int
}

// Registry holds all registered port groups.
type Registry struct {
	groups map[string]Group
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{groups: make(map[string]Group)}
}

// Add registers a group. Returns an error if the name is empty, the port list
// is empty, or a group with the same name already exists.
func (r *Registry) Add(g Group) error {
	if g.Name == "" {
		return errors.New("grouping: group name must not be empty")
	}
	if len(g.Ports) == 0 {
		return fmt.Errorf("grouping: group %q has no ports", g.Name)
	}
	if _, exists := r.groups[g.Name]; exists {
		return fmt.Errorf("grouping: group %q already registered", g.Name)
	}
	r.groups[g.Name] = g
	return nil
}

// Get returns the Group for the given name and whether it was found.
func (r *Registry) Get(name string) (Group, bool) {
	g, ok := r.groups[name]
	return g, ok
}

// GroupsForPort returns the names of all groups that contain the given port.
func (r *Registry) GroupsForPort(port int) []string {
	var names []string
	for name, g := range r.groups {
		for _, p := range g.Ports {
			if p == port {
				names = append(names, name)
				break
			}
		}
	}
	return names
}

// All returns every registered Group.
func (r *Registry) All() []Group {
	out := make([]Group, 0, len(r.groups))
	for _, g := range r.groups {
		out = append(out, g)
	}
	return out
}
