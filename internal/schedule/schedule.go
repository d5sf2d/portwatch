package schedule

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Runner periodically scans ports and alerts on changes.
type Runner struct {
	scanner  *scanner.Scanner
	store    *state.Store
	notifier *alert.Notifier
	ports    []int
	interval time.Duration
}

// New creates a new Runner.
func New(sc *scanner.Scanner, st *state.Store, n *alert.Notifier, ports []int, interval time.Duration) *Runner {
	return &Runner{
		scanner:  sc,
		store:    st,
		notifier: n,
		ports:    ports,
		interval: interval,
	}
}

// Run starts the periodic scan loop, blocking until ctx is cancelled.
func (r *Runner) Run(ctx context.Context) error {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	if err := r.tick(); err != nil {
		log.Printf("initial scan error: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := r.tick(); err != nil {
				log.Printf("scan error: %v", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *Runner) tick() error {
	prev, _ := r.store.Load()

	current, err := r.scanner.Scan(r.ports)
	if err != nil {
		return err
	}

	if prev != nil {
		diffs := state.Diff(prev, current)
		if len(diffs) > 0 {
			r.notifier.Notify(diffs)
		}
	}

	return r.store.Save(current)
}
