package portwatch

import (
	"io"
	"os"
	"time"
)

// Option is a functional option for Watcher.
type Option func(*Watcher)

// WithTimeout sets the per-port dial timeout used by the scanner.
func WithTimeout(d time.Duration) Option {
	return func(w *Watcher) {
		if d > 0 {
			w.scanner = newScannerWithTimeout(d)
		}
	}
}

// WithOutput redirects alert output to the given writer.
func WithOutput(out io.Writer) Option {
	return func(w *Watcher) {
		if out != nil {
			w.out = out
		}
	}
}

// WithLogPath sets the history log file path.
func WithLogPath(path string) Option {
	return func(w *Watcher) {
		w.applyLogPath(path)
	}
}

// applyLogPath is a helper so the method is testable independently.
func (w *Watcher) applyLogPath(path string) {
	if path == "" {
		return
	}
	importedLog, err := newHistoryLog(path)
	if err == nil {
		w.log = importedLog
	}
}

// defaultOutput returns os.Stdout if w is nil or w.out is nil.
func defaultOutput(w io.Writer) io.Writer {
	if w == nil {
		return os.Stdout
	}
	return w
}
