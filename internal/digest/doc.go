// Package digest provides lightweight SHA-256 fingerprinting of port
// snapshots. It is used by the scheduler to skip diff computation and
// alerting when consecutive scans produce identical results, reducing
// unnecessary processing on quiet hosts.
package digest
