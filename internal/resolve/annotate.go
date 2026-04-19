package resolve

import "github.com/user/portwatch/internal/state"

// AnnotatedPort pairs a port number with its resolved service name.
type AnnotatedPort struct {
	Port    int    `json:"port"`
	Service string `json:"service"`
	Known   bool   `json:"known"`
}

// AnnotateSnapshot returns annotated ports for every open port in the snapshot.
func AnnotateSnapshot(snap state.Snapshot, r *Resolver) []AnnotatedPort {
	out := make([]AnnotatedPort, 0, len(snap.OpenPorts))
	for _, p := range snap.OpenPorts {
		out = append(out, AnnotatedPort{
			Port:    p,
			Service: r.Name(p),
			Known:   r.Known(p),
		})
	}
	return out
}
