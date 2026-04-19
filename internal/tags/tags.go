// Package tags provides port tagging and label management for portwatch.
package tags

import "fmt"

// Tag associates a label and optional description with a port number.
type Tag struct {
	Port        int    `json:"port"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
}

// Registry holds a mapping from port numbers to tags.
type Registry struct {
	entries map[int]Tag
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{entries: make(map[int]Tag)}
}

// Add registers a tag for a port. Returns an error if the port is already tagged.
func (r *Registry) Add(t Tag) error {
	if t.Label == "" {
		return fmt.Errorf("tag label must not be empty for port %d", t.Port)
	}
	if _, exists := r.entries[t.Port]; exists {
		return fmt.Errorf("port %d already has a tag", t.Port)
	}
	r.entries[t.Port] = t
	return nil
}

// Get returns the Tag for a port and whether it was found.
func (r *Registry) Get(port int) (Tag, bool) {
	t, ok := r.entries[port]
	return t, ok
}

// Label returns the label for a port, or a default string if not tagged.
func (r *Registry) Label(port int) string {
	if t, ok := r.entries[port]; ok {
		return t.Label
	}
	return fmt.Sprintf("port-%d", port)
}

// All returns all registered tags.
func (r *Registry) All() []Tag {
	out := make([]Tag, 0, len(r.entries))
	for _, t := range r.entries {
		out = append(out, t)
	}
	return out
}

// Remove deletes the tag for a port. Returns an error if not found.
func (r *Registry) Remove(port int) error {
	if _, ok := r.entries[port]; !ok {
		return fmt.Errorf("no tag found for port %d", port)
	}
	delete(r.entries, port)
	return nil
}
