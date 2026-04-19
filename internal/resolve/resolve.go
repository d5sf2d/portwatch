// Package resolve maps port numbers to well-known service names.
package resolve

import "fmt"

// well-known port-to-service mapping.
var builtins = map[int]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// Resolver resolves port numbers to service names.
type Resolver struct {
	overrides map[int]string
}

// New returns a Resolver with optional user-supplied overrides.
func New(overrides map[int]string) *Resolver {
	if overrides == nil {
		overrides = map[int]string{}
	}
	return &Resolver{overrides: overrides}
}

// Name returns a human-readable service name for the given port.
// User overrides take precedence over builtins. Falls back to "port/<n>".
func (r *Resolver) Name(port int) string {
	if s, ok := r.overrides[port]; ok {
		return s
	}
	if s, ok := builtins[port]; ok {
		return s
	}
	return fmt.Sprintf("port/%d", port)
}

// Known reports whether the port maps to a named service.
func (r *Resolver) Known(port int) bool {
	_, inOverride := r.overrides[port]
	_, inBuiltin := builtins[port]
	return inOverride || inBuiltin
}
