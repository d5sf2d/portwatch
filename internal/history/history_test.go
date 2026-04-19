package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/portwatch/internal/history"
)

func tempLog(t *testing.T) *history.Log {
	t.Helper()
	dir := t.TempDir()
	return history.NewLog(filepath.Join(dir, "history.jsonl"))
}

func TestAppendAndReadAll(t *testing.T) {
	l := tempLog(t)
	now := time.Now().UTC().Truncate(time.Second)
	entries := []history.Entry{
		{Timestamp: now, Port: 80, Proto: "tcp", Event: "opened", Host: "localhost"},
		{Timestamp: now, Port: 443, Proto: "tcp", Event: "opened", Host: "localhost"},
	}
	if err := l.Append(entries); err != nil {
		t.Fatalf("Append: %v", err)
	}
	got, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 entries, got %d", len(got))
	}
	if got[0].Port != 80 {
		t.Errorf("want port 80, got %d", got[0].Port)
	}
}

func TestReadAll_MissingFile(t *testing.T) {
	l := history.NewLog(filepath.Join(t.TempDir(), "missing.jsonl"))
	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice for missing file")
	}
}

func TestAppend_Idempotent(t *testing.T) {
	l := tempLog(t)
	now := time.Now().UTC()
	e := history.Entry{Timestamp: now, Port: 22, Proto: "tcp", Event: "closed", Host: "localhost"}
	_ = l.Append([]history.Entry{e})
	_ = l.Append([]history.Entry{e})
	got, _ := l.ReadAll()
	if len(got) != 2 {
		t.Errorf("want 2 (append twice), got %d", len(got))
	}
}

func TestFilter(t *testing.T) {
	_ = os.Getenv // silence import
	now := time.Now().UTC()
	entries := []history.Entry{
		{Timestamp: now.Add(-2 * time.Hour), Port: 80, Event: "opened"},
		{Timestamp: now.Add(-1 * time.Hour), Port: 22, Event: "closed"},
		{Timestamp: now, Port: 80, Event: "closed"},
	}
	got := history.Filter(entries, history.FilterOptions{Port: 80})
	if len(got) != 2 {
		t.Errorf("want 2 port-80 entries, got %d", len(got))
	}
	got = history.Filter(entries, history.FilterOptions{Event: "closed"})
	if len(got) != 2 {
		t.Errorf("want 2 closed entries, got %d", len(got))
	}
	got = history.Filter(entries, history.FilterOptions{Since: now.Add(-90 * time.Minute)})
	if len(got) != 2 {
		t.Errorf("want 2 recent entries, got %d", len(got))
	}
}
