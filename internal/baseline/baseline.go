// Package baseline provides functionality to capture and compare
// a trusted baseline snapshot of open ports against current state.
package baseline

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// Baseline represents a trusted snapshot of ports at a point in time.
type Baseline struct {
	CapturedAt time.Time `json:"captured_at"`
	Ports      []int     `json:"ports"`
}

// Manager handles saving and loading baselines from disk.
type Manager struct {
	path string
}

// New returns a new Manager using the given file path.
func New(path string) *Manager {
	return &Manager{path: path}
}

// Save writes the baseline to disk.
func (m *Manager) Save(b Baseline) error {
	f, err := os.Create(m.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(b)
}

// Load reads the baseline from disk.
func (m *Manager) Load() (Baseline, error) {
	f, err := os.Open(m.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Baseline{}, ErrNoBaseline
		}
		return Baseline{}, err
	}
	defer f.Close()
	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return Baseline{}, err
	}
	return b, nil
}

// Capture creates a new Baseline from the given list of open ports.
func Capture(ports []int) Baseline {
	cp := make([]int, len(ports))
	copy(cp, ports)
	return Baseline{
		CapturedAt: time.Now().UTC(),
		Ports:      cp,
	}
}

// ErrNoBaseline is returned when no baseline file exists.
var ErrNoBaseline = errors.New("baseline: no baseline file found")
