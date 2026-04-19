// Package audit provides a structured audit trail for port scan events,
// recording who triggered a scan, when, and what changed.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Trigger   string    `json:"trigger"`
	PortsScanned []int  `json:"ports_scanned"`
	ChangesFound int    `json:"changes_found"`
	Note      string    `json:"note,omitempty"`
}

// Log writes audit entries to a newline-delimited JSON file.
type Log struct {
	mu   sync.Mutex
	path string
}

// NewLog creates a new audit Log backed by the file at path.
func NewLog(path string) *Log {
	return &Log{path: path}
}

// Record appends an entry to the audit log.
func (l *Log) Record(e Entry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}

	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(e); err != nil {
		return fmt.Errorf("audit: encode entry: %w", err)
	}
	return nil
}

// ReadAll returns all audit entries from the log file.
func (l *Log) ReadAll() ([]Entry, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	f, err := os.Open(l.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return fmt.Errorf("audit: open file: %w", err)
	}
	defer f.Close()

	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
