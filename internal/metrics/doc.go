// Package metrics provides an in-memory counter collector for portwatch
// scan-cycle telemetry. It tracks total scans, alerts fired, ports observed
// open, and suppressed notifications, and exposes them via a thread-safe
// Snapshot for use in status reports or health endpoints.
package metrics
