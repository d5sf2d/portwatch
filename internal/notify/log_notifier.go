package notify

import (
	"fmt"
	"io"
	"os"
	"time"
)

// LogNotifier writes notifications as plain-text lines to a writer.
type LogNotifier struct {
	w io.Writer
}

// NewLogNotifier returns a LogNotifier writing to w.
// If w is nil, os.Stdout is used.
func NewLogNotifier(w io.Writer) *LogNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &LogNotifier{w: w}
}

// Send formats msg and writes it to the underlying writer.
func (l *LogNotifier) Send(msg Message) error {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	_, err := fmt.Fprintf(l.w, "[%s] [%s] %s: %s\n",
		timestamp, msg.Level, msg.Title, msg.Body)
	return err
}

// Name returns the backend identifier.
func (l *LogNotifier) Name() string { return "log" }
