// Package debounce prevents alert flooding by suppressing repeated identical
// diffs within a configurable time window.
package debounce

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Debouncer tracks recently seen diffs and suppresses duplicates within a window.
type Debouncer struct {
	mu      sync.Mutex
	seen    map[string]time.Time
	window  time.Duration
	nowFunc func() time.Time
}

// New returns a Debouncer with the given suppression window.
func New(window time.Duration) *Debouncer {
	return &Debouncer{
		seen:    make(map[string]time.Time),
		window:  window,
		nowFunc: time.Now,
	}
}

// key builds a stable string key for a diff entry.
func key(d state.Diff) string {
	return fmt.Sprintf("%s:%d:%s", d.Host, d.Port, d.Change)
}

// IsDuplicate reports whether the diff has already been seen within the window.
// If it has not been seen (or the window has expired), it records the diff and
// returns false so the caller proceeds with alerting.
func (db *Debouncer) IsDuplicate(d state.Diff) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	now := db.nowFunc()
	k := key(d)

	if last, ok := db.seen[k]; ok && now.Sub(last) < db.window {
		return true
	}

	db.seen[k] = now
	return false
}

// Filter returns only the diffs that are not duplicates within the window.
func (db *Debouncer) Filter(diffs []state.Diff) []state.Diff {
	out := diffs[:0:0]
	for _, d := range diffs {
		if !db.IsDuplicate(d) {
			out = append(out, d)
		}
	}
	return out
}

// Expire removes all entries whose window has elapsed.
func (db *Debouncer) Expire() {
	db.mu.Lock()
	defer db.mu.Unlock()

	now := db.nowFunc()
	for k, t := range db.seen {
		if now.Sub(t) >= db.window {
			delete(db.seen, k)
		}
	}
}
