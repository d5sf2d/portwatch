package state

import (
	"os"
	"testing"
	"time"
)

func makeSnapshot(host string, ports []PortState) Snapshot {
	return Snapshot{Timestamp: time.Now(), Host: host, Ports: ports}
}

func TestStore_SaveLoad(t *testing.T) {
	f, err := os.CreateTemp("", "portwatch-state-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	store := NewStore(f.Name())
	snap := makeSnapshot("localhost", []PortState{
		{Port: 80, Protocol: "tcp", Open: true},
		{Port: 443, Protocol: "tcp", Open: true},
	})

	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Host != snap.Host {
		t.Errorf("host mismatch: got %q want %q", loaded.Host, snap.Host)
	}
	if len(loaded.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(loaded.Ports))
	}
}

func TestStore_LoadMissing(t *testing.T) {
	store := NewStore("/tmp/portwatch-nonexistent-state.json")
	snap, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap != nil {
		t.Error("expected nil snapshot for missing file")
	}
}

func TestDiff_NewOpenPort(t *testing.T) {
	prev := makeSnapshot("localhost", []PortState{})
	curr := makeSnapshot("localhost", []PortState{{Port: 8080, Protocol: "tcp", Open: true}})
	changes := Diff(prev, curr)
	if len(changes) != 1 || !changes[0].Opened || changes[0].Port != 8080 {
		t.Errorf("expected one Opened change for port 8080, got %+v", changes)
	}
}

func TestDiff_PortClosed(t *testing.T) {
	prev := makeSnapshot("localhost", []PortState{{Port: 22, Protocol: "tcp", Open: true}})
	curr := makeSnapshot("localhost", []PortState{{Port: 22, Protocol: "tcp", Open: false}})
	changes := Diff(prev, curr)
	if len(changes) != 1 || !changes[0].Closed || changes[0].Port != 22 {
		t.Errorf("expected one Closed change for port 22, got %+v", changes)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	ports := []PortState{{Port: 443, Protocol: "tcp", Open: true}}
	prev := makeSnapshot("localhost", ports)
	curr := makeSnapshot("localhost", ports)
	if changes := Diff(prev, curr); len(changes) != 0 {
		t.Errorf("expected no changes, got %+v", changes)
	}
}
