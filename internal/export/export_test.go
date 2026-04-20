package export_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/export"
)

func makeSnapshot() export.Snapshot {
	return export.Snapshot{
		Host:      "localhost",
		ScannedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Ports: []export.PortEntry{
			{Port: 22, Proto: "tcp", Service: "ssh", Tag: "infra"},
			{Port: 80, Proto: "tcp", Service: "http"},
		},
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	e := export.New(&buf, export.FormatJSON)
	if err := e.Write(makeSnapshot()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var snap export.Snapshot
	if err := json.Unmarshal(buf.Bytes(), &snap); err != nil {
		t.Fatalf("invalid json output: %v", err)
	}
	if snap.Host != "localhost" {
		t.Errorf("expected host localhost, got %q", snap.Host)
	}
	if len(snap.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(snap.Ports))
	}
}

func TestWrite_CSVFormat(t *testing.T) {
	var buf bytes.Buffer
	e := export.New(&buf, export.FormatCSV)
	if err := e.Write(makeSnapshot()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + 2 data rows
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(lines), lines)
	}
	if !strings.HasPrefix(lines[0], "host,") {
		t.Errorf("expected csv header, got %q", lines[0])
	}
	if !strings.Contains(lines[1], "22") {
		t.Errorf("expected port 22 in row 1, got %q", lines[1])
	}
	if !strings.Contains(lines[2], "80") {
		t.Errorf("expected port 80 in row 2, got %q", lines[2])
	}
}

func TestWrite_NilWriterDiscard(t *testing.T) {
	e := export.New(nil, export.FormatJSON)
	if err := e.Write(makeSnapshot()); err != nil {
		t.Fatalf("unexpected error with nil writer: %v", err)
	}
}

func TestWrite_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	e := export.New(&buf, export.Format("xml"))
	if err := e.Write(makeSnapshot()); err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestWrite_EmptyPorts(t *testing.T) {
	var buf bytes.Buffer
	e := export.New(&buf, export.FormatCSV)
	snap := export.Snapshot{Host: "host1", ScannedAt: time.Now(), Ports: []export.PortEntry{}}
	if err := e.Write(snap); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header line, got %d lines", len(lines))
	}
}
