// Package notify provides pluggable notification backends for portwatch.
package notify

import "fmt"

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Message holds the data sent to a notifier.
type Message struct {
	Level   Level
	Title   string
	Body    string
	Tags    []string
}

// Notifier is the interface implemented by all notification backends.
type Notifier interface {
	Send(msg Message) error
	Name() string
}

// Multi fans out a message to multiple notifiers, collecting errors.
type Multi struct {
	notifiers []Notifier
}

// NewMulti returns a Multi notifier wrapping the provided backends.
func NewMulti(nn ...Notifier) *Multi {
	return &Multi{notifiers: nn}
}

// Send delivers msg to every registered notifier.
func (m *Multi) Send(msg Message) error {
	var errs []error
	for _, n := range m.notifiers {
		if err := n.Send(msg); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", n.Name(), err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("notify errors: %v", errs)
	}
	return nil
}

// Name returns the identifier for the multi-notifier.
func (m *Multi) Name() string { return "multi" }
