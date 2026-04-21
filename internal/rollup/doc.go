// Package rollup provides a time-windowed diff aggregator for portwatch.
//
// Instead of emitting one alert per port change, callers can feed diffs
// into a Rollup and receive a batched Summary once the window closes or
// the buffer is full. This reduces notification noise during large
// network-reconfiguration events.
package rollup
