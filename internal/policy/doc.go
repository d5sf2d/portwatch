// Package policy provides rule-based evaluation of open ports.
//
// Rules are defined as JSON and loaded via LoadFile. Each rule can
// target a specific port, a set of hosts, or apply globally. When a
// scan result violates a rule an Evaluator returns Violation values
// that the caller can forward to the alert or notify subsystems.
package policy
