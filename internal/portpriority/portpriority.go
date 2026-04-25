// Package portpriority assigns scan and alert priority levels to ports
// based on their classification and known service criticality.
package portpriority

import "fmt"

// Level represents a priority level for a port.
type Level int

const (
	Low    Level = 1
	Medium Level = 2
	High   Level = 3
	Critical Level = 4
)

func (l Level) String() string {
	switch l {
	case Low:
		return "low"
	case Medium:
		return "medium"
	case High:
		return "high"
	case Critical:
		return "critical"
	default:
		return fmt.Sprintf("unknown(%d)", int(l))
	}
}

// Registry maps ports to priority levels.
type Registry struct {
	overrides map[int]Level
}

// New returns a new Registry with built-in defaults.
func New() *Registry {
	return &Registry{
		overrides: make(map[int]Level),
	}
}

// Set overrides the priority for a specific port.
func (r *Registry) Set(port int, level Level) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portpriority: invalid port %d", port)
	}
	r.overrides[port] = level
	return nil
}

// Get returns the priority level for the given port.
// Overrides take precedence over built-in defaults.
func (r *Registry) Get(port int) Level {
	if lvl, ok := r.overrides[port]; ok {
		return lvl
	}
	return defaultLevel(port)
}

// defaultLevel derives a priority from well-known port ranges and services.
func defaultLevel(port int) Level {
	switch port {
	case 22, 23, 3389: // SSH, Telnet, RDP
		return Critical
	case 21, 25, 110, 143, 993, 995: // FTP, SMTP, POP3, IMAP
		return High
	case 80, 443, 8080, 8443: // HTTP/S
		return High
	case 3306, 5432, 1433, 27017, 6379: // DB ports
		return Critical
	}
	switch {
	case port < 1024:
		return Medium
	case port < 49152:
		return Low
	default:
		return Low
	}
}
