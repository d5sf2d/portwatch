// Package policy evaluates scan results against user-defined rules
// and produces a list of violations for alerting or reporting.
package policy

import (
	"fmt"
	"time"
)

// Rule describes a single policy constraint.
type Rule struct {
	Name        string // human-readable name
	Port        int    // 0 means any port
	MustBeClosed bool   // true → port must not be open
	AllowedHosts []string // empty means all hosts
}

// Violation is produced when a scan result breaks a rule.
type Violation struct {
	Rule      Rule
	Port      int
	Host      string
	DetectedAt time.Time
}

func (v Violation) String() string {
	return fmt.Sprintf("policy %q violated: port %d on %s is open", v.Rule.Name, v.Port, v.Host)
}

// Evaluator holds a set of rules and checks snapshots against them.
type Evaluator struct {
	rules []Rule
}

// New returns an Evaluator loaded with the provided rules.
func New(rules []Rule) *Evaluator {
	return &Evaluator{rules: rules}
}

// Evaluate checks each open port in ports (host → []port) against every rule
// and returns any violations found.
func (e *Evaluator) Evaluate(host string, ports []int) []Violation {
	now := time.Now()
	var violations []Violation

	for _, rule := range e.rules {
		if !rule.appliesToHost(host) {
			continue
		}
		for _, p := range ports {
			if rule.Port != 0 && rule.Port != p {
				continue
			}
			if rule.MustBeClosed {
				violations = append(violations, Violation{
					Rule:       rule,
					Port:       p,
					Host:       host,
					DetectedAt: now,
				})
			}
		}
	}
	return violations
}

func (r Rule) appliesToHost(host string) bool {
	if len(r.AllowedHosts) == 0 {
		return true
	}
	for _, h := range r.AllowedHosts {
		if h == host {
			return true
		}
	}
	return false
}
