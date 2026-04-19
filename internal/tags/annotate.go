package tags

import "github.com/user/portwatch/internal/state"

// AnnotatedPort pairs a port number with its resolved label.
type AnnotatedPort struct {
	Port  int
	Label string
	Open  bool
}

// Annotate enriches a state.Snapshot's ports with labels from the registry.
func Annotate(snap state.Snapshot, reg *Registry) []AnnotatedPort {
	out := make([]AnnotatedPort, 0, len(snap.Ports))
	for _, p := range snap.Ports {
		out = append(out, AnnotatedPort{
			Port:  p.Port,
			Label: reg.Label(p.Port),
			Open:  p.Open,
		})
	}
	return out
}
