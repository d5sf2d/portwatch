// Package portlabel provides human-readable severity labels for ports
// based on their classification and known service associations.
package portlabel

import "fmt"

// Severity represents the alert severity level for a port.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

// Label holds the human-readable metadata for a port.
type Label struct {
	Port        int
	Service     string
	Severity    Severity
	Description string
}

// String returns a formatted representation of the label.
func (l Label) String() string {
	return fmt.Sprintf("port=%d service=%s severity=%s", l.Port, l.Service, l.Severity)
}

// Labeler assigns severity labels to ports.
type Labeler struct {
	overrides map[int]Label
}

// New returns a Labeler with built-in defaults.
func New() *Labeler {
	return &Labeler{
		overrides: make(map[int]Label),
	}
}

// Add registers a custom label for a port, overriding built-in defaults.
func (l *Labeler) Add(port int, service string, sev Severity, desc string) {
	l.overrides[port] = Label{Port: port, Service: service, Severity: sev, Description: desc}
}

// Label returns the Label for a given port.
func (l *Labeler) Label(port int) Label {
	if lbl, ok := l.overrides[port]; ok {
		return lbl
	}
	return defaultLabel(port)
}

// defaultLabel returns a built-in label based on well-known port ranges.
func defaultLabel(port int) Label {
	switch port {
	case 22:
		return Label{Port: port, Service: "ssh", Severity: SeverityHigh, Description: "Secure Shell"}
	case 23:
		return Label{Port: port, Service: "telnet", Severity: SeverityCritical, Description: "Telnet (unencrypted)"}
	case 80:
		return Label{Port: port, Service: "http", Severity: SeverityMedium, Description: "HTTP"}
	case 443:
		return Label{Port: port, Service: "https", Severity: SeverityLow, Description: "HTTPS"}
	case 3306:
		return Label{Port: port, Service: "mysql", Severity: SeverityCritical, Description: "MySQL database"}
	case 5432:
		return Label{Port: port, Service: "postgres", Severity: SeverityCritical, Description: "PostgreSQL database"}
	case 6379:
		return Label{Port: port, Service: "redis", Severity: SeverityHigh, Description: "Redis"}
	}
	switch {
	case port < 1024:
		return Label{Port: port, Service: "unknown", Severity: SeverityMedium, Description: "System port"}
	case port < 49152:
		return Label{Port: port, Service: "unknown", Severity: SeverityLow, Description: "Registered port"}
	default:
		return Label{Port: port, Service: "unknown", Severity: SeverityInfo, Description: "Dynamic port"}
	}
}
