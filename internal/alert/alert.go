package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event describes a single alert event derived from a port diff.
type Event struct {
	Level     Level
	Message   string
	Timestamp time.Time
}

// Notifier writes alert events to an output destination.
type Notifier struct {
	out io.Writer
}

// New returns a Notifier that writes to w. Pass nil to use os.Stdout.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify converts a slice of Diff entries into Events and writes them.
func (n *Notifier) Notify(diffs []state.DiffEntry) []Event {
	var events []Event
	for _, d := range diffs {
		var ev Event
		ev.Timestamp = time.Now()
		switch d.Kind {
		case state.DiffOpened:
			ev.Level = LevelAlert
			ev.Message = fmt.Sprintf("port %d/%s opened (service: %s)", d.Port, d.Proto, d.Service)
		case state.DiffClosed:
			ev.Level = LevelWarn
			ev.Message = fmt.Sprintf("port %d/%s closed (service: %s)", d.Port, d.Proto, d.Service)
		default:
			continue
		}
		events = append(events, ev)
		fmt.Fprintf(n.out, "[%s] %s %s\n", ev.Timestamp.Format(time.RFC3339), ev.Level, ev.Message)
	}
	return events
}
