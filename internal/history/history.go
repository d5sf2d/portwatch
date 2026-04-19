// Package history maintains a rolling log of scan events for trend analysis.
package history

import (
	"encoding/json"
	"os"
	"time"
)

// Entry records a single scan event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Port      int       `json:"port"`
	Proto     string    `json:"proto"`
	Event     string    `json:"event"` // "opened" | "closed"
	Host      string    `json:"host"`
}

// Log is an append-only history log backed by a JSON file.
type Log struct {
	path string
}

// NewLog returns a Log that persists entries to path.
func NewLog(path string) *Log {
	return &Log{path: path}
}

// Append adds entries to the log file.
func (l *Log) Append(entries []Entry) error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, e := range entries {
		if err := enc.Encode(e); err != nil {
			return err
		}
	}
	return nil
}

// ReadAll returns all entries stored in the log file.
func (l *Log) ReadAll() ([]Entry, error) {
	f, err := os.Open(l.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}
