package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempLog(t *testing.T) (*Log, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")
	return NewLog(path), path
}

func TestRecord_And_ReadAll(t *testing.T) {
	log, _ := tempLog(t)

	e := Entry{
		Timestamp:    time.Now().UTC(),
		Trigger:      "scheduler",
		PortsScanned: []int{80, 443},
		ChangesFound: 1,
	}
	if err := log.Record(e); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := log.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Trigger != "scheduler" {
		t.Errorf("trigger mismatch: got %q", entries[0].Trigger)
	}
	if entries[0].ChangesFound != 1 {
		t.Errorf("changes mismatch: got %d", entries[0].ChangesFound)
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	log, _ := tempLog(t)

	for i := 0; i < 3; i++ {
		if err := log.Record(Entry{Trigger: "manual", PortsScanned: []int{22}}); err != nil {
			t.Fatalf("Record %d: %v", i, err)
		}
	}

	entries, err := log.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestReadAll_MissingFile(t *testing.T) {
	log := NewLog("/tmp/portwatch_audit_nonexistent_xyz.jsonl")
	defer os.Remove("/tmp/portwatch_audit_nonexistent_xyz.jsonl")

	entries, err := log.ReadAll()
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries, got %v", entries)
	}
}

func TestRecord_SetsTimestampIfZero(t *testing.T) {
	log, _ := tempLog(t)

	if err := log.Record(Entry{Trigger: "test"}); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, _ := log.ReadAll()
	if entries[0].Timestamp.IsZero() {
		t.Error("expected timestamp to be set automatically")
	}
}
