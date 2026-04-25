package portwatch

// internal.go contains thin constructor wrappers so that option.go and
// portwatch.go can reference external packages without import cycles and
// remain easy to stub in tests.

import (
	"time"

	"github.com/example/portwatch/internal/history"
	"github.com/example/portwatch/internal/scanner"
)

// newScannerWithTimeout returns a scanner.Scanner configured with d.
func newScannerWithTimeout(d time.Duration) *scanner.Scanner {
	return scanner.New(d)
}

// newHistoryLog wraps history.NewLog so callers don't need the import.
func newHistoryLog(path string) (*history.Log, error) {
	return history.NewLog(path)
}
