package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/state"
)

func TestNotify_OpenedPort(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	diffs := []state.DiffEntry{
		{Kind: state.DiffOpened, Port: 8080, Proto: "tcp", Service: "http-alt"},
	}

	events := n.Notify(diffs)

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Level != alert.LevelAlert {
		t.Errorf("expected level ALERT, got %s", events[0].Level)
	}
	if !strings.Contains(buf.String(), "8080") {
		t.Errorf("expected port 8080 in output, got: %s", buf.String())
	}
}

func TestNotify_ClosedPort(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	diffs := []state.DiffEntry{
		{Kind: state.DiffClosed, Port: 22, Proto: "tcp", Service: "ssh"},
	}

	events := n.Notify(diffs)

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Level != alert.LevelWarn {
		t.Errorf("expected level WARN, got %s", events[0].Level)
	}
	if !strings.Contains(buf.String(), "ssh") {
		t.Errorf("expected service name in output, got: %s", buf.String())
	}
}

func TestNotify_NilWriterDefaultsToStdout(t *testing.T) {
	// Should not panic when nil writer is passed.
	n := alert.New(nil)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}

func TestNotify_EmptyDiffs(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)
	events := n.Notify(nil)
	if len(events) != 0 {
		t.Errorf("expected 0 events for empty diffs, got %d", len(events))
	}
}
