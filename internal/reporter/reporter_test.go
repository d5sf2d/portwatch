package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/state"
)

func makeSnap() state.Snapshot {
	return state.Snapshot{
		Host:      "localhost",
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Ports: []state.PortInfo{
			{Port: 80, Service: "http"},
			{Port: 443, Service: "https"},
		},
	}
}

func TestWrite_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Write(makeSnap()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "localhost") {
		t.Error("expected host in output")
	}
	if !strings.Contains(out, "2 open port") {
		t.Error("expected port count in output")
	}
	if !strings.Contains(out, "http") {
		t.Error("expected service name in output")
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatJSON)
	if err := r.Write(makeSnap()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var rep reporter.Report
	if err := json.Unmarshal(buf.Bytes(), &rep); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if rep.Total != 2 {
		t.Errorf("expected total 2, got %d", rep.Total)
	}
	if rep.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", rep.Host)
	}
}

func TestNew_NilWriterDefaultsToStdout(t *testing.T) {
	// Should not panic
	r := reporter.New(nil, reporter.FormatText)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}

func TestWrite_EmptyPorts(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	snap := state.Snapshot{Host: "localhost", Timestamp: time.Now(), Ports: []state.PortInfo{}}
	if err := r.Write(snap); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "0 open port") {
		t.Error("expected zero port count")
	}
}
