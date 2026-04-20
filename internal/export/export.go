// Package export provides functionality for exporting port scan snapshots
// to various file formats for offline analysis or integration with other tools.
package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Snapshot represents a point-in-time view of scanned ports.
type Snapshot struct {
	Host      string    `json:"host"`
	ScannedAt time.Time `json:"scanned_at"`
	Ports     []PortEntry `json:"ports"`
}

// PortEntry holds metadata for a single open port.
type PortEntry struct {
	Port    int    `json:"port"`
	Proto   string `json:"proto"`
	Service string `json:"service,omitempty"`
	Tag     string `json:"tag,omitempty"`
}

// Format defines the output format for an export.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Exporter writes snapshots to an io.Writer in a specified format.
type Exporter struct {
	w      io.Writer
	format Format
}

// New creates a new Exporter. If w is nil, output is discarded.
func New(w io.Writer, format Format) *Exporter {
	if w == nil {
		w = io.Discard
	}
	return &Exporter{w: w, format: format}
}

// Write serializes the snapshot to the configured format.
func (e *Exporter) Write(snap Snapshot) error {
	switch e.format {
	case FormatJSON:
		return e.writeJSON(snap)
	case FormatCSV:
		return e.writeCSV(snap)
	default:
		return fmt.Errorf("export: unsupported format %q", e.format)
	}
}

func (e *Exporter) writeJSON(snap Snapshot) error {
	enc := json.NewEncoder(e.w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("export: json encode: %w", err)
	}
	return nil
}

func (e *Exporter) writeCSV(snap Snapshot) error {
	w := csv.NewWriter(e.w)
	if err := w.Write([]string{"host", "scanned_at", "port", "proto", "service", "tag"}); err != nil {
		return fmt.Errorf("export: csv header: %w", err)
	}
	ts := snap.ScannedAt.Format(time.RFC3339)
	for _, p := range snap.Ports {
		row := []string{
			snap.Host,
			ts,
			fmt.Sprintf("%d", p.Port),
			p.Proto,
			p.Service,
			p.Tag,
		}
		if err := w.Write(row); err != nil {
			return fmt.Errorf("export: csv row: %w", err)
		}
	}
	w.Flush()
	return w.Error()
}
