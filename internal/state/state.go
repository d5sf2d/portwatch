package state

import (
	"encoding/json"
	"os"
	"time"
)

// PortState represents the recorded state of a single port.
type PortState struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Open     bool   `json:"open"`
}

// Snapshot holds a full scan result at a point in time.
type Snapshot struct {
	Timestamp time.Time   `json:"timestamp"`
	Host      string      `json:"host"`
	Ports     []PortState `json:"ports"`
}

// Store persists and loads snapshots to/from a JSON file.
type Store struct {
	path string
}

// NewStore creates a Store backed by the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the snapshot to disk, overwriting any previous state.
func (s *Store) Save(snap Snapshot) error {
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}

// Load reads the last snapshot from disk.
// Returns nil, nil when no previous state exists.
func (s *Store) Load() (*Snapshot, error) {
	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, err
	}
	return &snap, nil
}
