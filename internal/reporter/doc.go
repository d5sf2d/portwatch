// Package reporter provides formatted output of port scan snapshots.
//
// It supports two output formats:
//   - text: human-readable, suitable for terminal display
//   - json: machine-readable, suitable for log ingestion or piping
//
// Usage:
//
//	r := reporter.New(os.Stdout, reporter.FormatText)
//	r.Write(snapshot)
package reporter
