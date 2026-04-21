package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Format controls the output format of the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds a snapshot summary for output.
type Report struct {
	Timestamp time.Time        `json:"timestamp"`
	Host      string           `json:"host"`
	OpenPorts []state.PortInfo `json:"open_ports"`
	Total     int              `json:"total"`
}

// Reporter writes port scan reports to a writer.
type Reporter struct {
	out    io.Writer
	format Format
}

// New creates a Reporter. If out is nil, os.Stdout is used.
func New(out io.Writer, format Format) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	if format == "" {
		format = FormatText
	}
	return &Reporter{out: out, format: format}
}

// Write outputs the snapshot as a report.
func (r *Reporter) Write(snap state.Snapshot) error {
	report := Report{
		Timestamp: snap.Timestamp,
		Host:      snap.Host,
		OpenPorts: snap.Ports,
		Total:     len(snap.Ports),
	}
	switch r.format {
	case FormatJSON:
		return r.writeJSON(report)
	default:
		return r.writeText(report)
	}
}

// WriteDiff outputs a human-readable summary of port changes between two snapshots.
func (r *Reporter) WriteDiff(diff state.Diff) error {
	if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
		return nil
	}
	_, err := fmt.Fprintf(r.out, "[%s] Host: %s — port changes detected\n",
		time.Now().Format(time.RFC3339), diff.Host)
	if err != nil {
		return err
	}
	for _, p := range diff.Opened {
		_, err = fmt.Fprintf(r.out, "  + %-6d %s\n", p.Port, p.Service)
		if err != nil {
			return err
		}
	}
	for _, p := range diff.Closed {
		_, err = fmt.Fprintf(r.out, "  - %-6d %s\n", p.Port, p.Service)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Reporter) writeText(rep Report) error {
	_, err := fmt.Fprintf(r.out, "[%s] Host: %s — %d open port(s)\n",
		rep.Timestamp.Format(time.RFC3339), rep.Host, rep.Total)
	if err != nil {
		return err
	}
	for _, p := range rep.OpenPorts {
		_, err = fmt.Fprintf(r.out, "  %-6d %s\n", p.Port, p.Service)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Reporter) writeJSON(rep Report) error {
	enc := json.NewEncoder(r.out)
	enc.SetIndent("", "  ")
	return enc.Encode(rep)
}
