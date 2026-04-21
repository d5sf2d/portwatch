// Package portmap provides a human-readable mapping from well-known port
// numbers to service names, and helpers to annotate scan results.
package portmap

// Registry maps port numbers to canonical service names.
type Registry struct {
	entries map[int]string
}

// wellKnown is the built-in set of common port-to-service mappings.
var wellKnown = map[int]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	465:  "smtps",
	587:  "submission",
	993:  "imaps",
	995:  "pop3s",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// New returns a Registry pre-loaded with well-known mappings.
func New() *Registry {
	entries := make(map[int]string, len(wellKnown))
	for k, v := range wellKnown {
		entries[k] = v
	}
	return &Registry{entries: entries}
}

// Add registers a custom port-to-name mapping, overriding any existing entry.
// An empty name is ignored.
func (r *Registry) Add(port int, name string) {
	if name == "" {
		return
	}
	r.entries[port] = name
}

// Lookup returns the service name for the given port and whether it was found.
func (r *Registry) Lookup(port int) (string, bool) {
	name, ok := r.entries[port]
	return name, ok
}

// LookupDefault returns the service name or the provided fallback string.
func (r *Registry) LookupDefault(port int, fallback string) string {
	if name, ok := r.entries[port]; ok {
		return name
	}
	return fallback
}

// All returns a copy of all current port mappings.
func (r *Registry) All() map[int]string {
	copy := make(map[int]string, len(r.entries))
	for k, v := range r.entries {
		copy[k] = v
	}
	return copy
}
