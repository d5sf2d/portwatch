// Package history provides an append-only event log for portwatch scan
// results. Each time the scheduler detects a port change, the caller may
// record an Entry so that operators can review historical open/close events
// and spot recurring patterns over time.
//
// Entries are stored as newline-delimited JSON so the file remains human-
// readable and can be processed with standard tools such as jq.
package history
