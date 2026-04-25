// Package portwatch provides a high-level watcher that ties together
// scanning, diffing, alerting, and history in a single reusable component.
package portwatch

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/history"
	"github.com/example/portwatch/internal/scanner"
	"github.com/example/portwatch/internal/snapshot"
	"github.com/example/portwatch/internal/state"
)

// Watcher performs periodic port scans and emits diffs.
type Watcher struct {
	host    string
	ports   []int
	store   *state.Store
	scanner *scanner.Scanner
	alerter *alert.Alerter
	log     *history.Log
	out     io.Writer
}

// Config holds watcher configuration.
type Config struct {
	Host     string
	Ports    []int
	StateDir string
	LogPath  string
	Out      io.Writer
}

// New constructs a Watcher from the given Config.
func New(cfg Config) (*Watcher, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("portwatch: host must not be empty")
	}
	if len(cfg.Ports) == 0 {
		return nil, fmt.Errorf("portwatch: at least one port required")
	}

	store := state.NewStore(cfg.StateDir)
	sc := scanner.New(2 * time.Second)
	al := alert.New(cfg.Out)

	var lg *history.Log
	if cfg.LogPath != "" {
		var err error
		lg, err = history.NewLog(cfg.LogPath)
		if err != nil {
			return nil, fmt.Errorf("portwatch: open history log: %w", err)
		}
	}

	return &Watcher{
		host:    cfg.Host,
		ports:   cfg.Ports,
		store:   store,
		scanner: sc,
		alerter: al,
		log:     lg,
		out:     cfg.Out,
	}, nil
}

// Run performs one scan cycle: scan → diff → alert → persist.
func (w *Watcher) Run(ctx context.Context) error {
	open, err := w.scanner.Scan(ctx, w.host, w.ports)
	if err != nil {
		return fmt.Errorf("portwatch: scan: %w", err)
	}

	curr := snapshot.New(w.host, open)

	prev, _ := w.store.Load(w.host)
	diffs := state.Diff(prev, curr)

	if len(diffs) > 0 {
		w.alerter.Notify(diffs)
		if w.log != nil {
			for _, d := range history.FromDiffs(w.host, diffs) {
				_ = w.log.Append(d)
			}
		}
	}

	return w.store.Save(curr)
}
