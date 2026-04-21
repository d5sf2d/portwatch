// Package pipeline wires together the core portwatch components into a
// single reusable scan-diff-alert execution unit. Callers construct a
// Pipeline with the desired dependencies and call Run to perform one
// complete cycle: scan → diff → suppress → debounce → alert → record.
package pipeline

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/audit"
	"github.com/user/portwatch/internal/debounce"
	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
	"github.com/user/portwatch/internal/suppress"
)

// Result summarises the outcome of a single pipeline run.
type Result struct {
	Host        string
	ScannedAt   time.Time
	OpenPorts   []int
	DiffCount   int
	Suppressed  int
	AlertsSent  int
}

// Pipeline holds the wired-up dependencies for one scan cycle.
type Pipeline struct {
	scanner  *scanner.Scanner
	store    *state.Store
	alerter  *alert.Alerter
	auditor  *audit.Log
	met      *metrics.Metrics
	debounce *debounce.Debouncer
	suppress *suppress.List
	host     string
	ports    []int
	out      io.Writer
}

// Config carries the options used to build a Pipeline.
type Config struct {
	Host     string
	Ports    []int
	Scanner  *scanner.Scanner
	Store    *state.Store
	Alerter  *alert.Alerter
	Auditor  *audit.Log
	Metrics  *metrics.Metrics
	Debounce *debounce.Debouncer
	Suppress *suppress.List
	// Out receives human-readable progress lines; defaults to io.Discard.
	Out io.Writer
}

// New validates cfg and returns a ready-to-use Pipeline.
func New(cfg Config) (*Pipeline, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("pipeline: host must not be empty")
	}
	if len(cfg.Ports) == 0 {
		return nil, fmt.Errorf("pipeline: at least one port is required")
	}
	if cfg.Scanner == nil {
		return nil, fmt.Errorf("pipeline: scanner is required")
	}
	if cfg.Store == nil {
		return nil, fmt.Errorf("pipeline: store is required")
	}
	out := cfg.Out
	if out == nil {
		out = io.Discard
	}
	return &Pipeline{
		scanner:  cfg.Scanner,
		store:    cfg.Store,
		alerter:  cfg.Alerter,
		auditor:  cfg.Auditor,
		met:      cfg.Metrics,
		debounce: cfg.Debounce,
		suppress: cfg.Suppress,
		host:     cfg.Host,
		ports:    cfg.Ports,
		out:      out,
	}, nil
}

// Run executes one full scan cycle and returns a Result.
// The provided context is forwarded to the scanner so callers can
// enforce a per-cycle deadline.
func (p *Pipeline) Run(ctx context.Context) (Result, error) {
	now := time.Now()
	res := Result{Host: p.host, ScannedAt: now}

	// 1. Scan
	snap, err := p.scanner.Scan(ctx, p.host, p.ports)
	if err != nil {
		return res, fmt.Errorf("pipeline: scan: %w", err)
	}
	for _, pr := range snap.Ports {
		if pr.Open {
			res.OpenPorts = append(res.OpenPorts, pr.Port)
		}
	}
	fmt.Fprintf(p.out, "[%s] scanned %d ports, %d open\n",
		now.Format(time.RFC3339), len(p.ports), len(res.OpenPorts))

	if p.met != nil {
		p.met.RecordScan(p.host)
	}

	// 2. Load previous snapshot and diff
	prev, _ := p.store.Load(p.host) // missing == first run, diff will show all as new
	diffs := state.Diff(prev, snap)
	res.DiffCount = len(diffs)

	// 3. Suppress maintenance windows
	if p.suppress != nil {
		before := len(diffs)
		diffs = suppress.FilterDiffs(diffs, p.suppress)
		res.Suppressed = before - len(diffs)
		if p.met != nil && res.Suppressed > 0 {
			p.met.RecordSuppressed(p.host, res.Suppressed)
		}
	}

	// 4. Debounce repeated flaps
	if p.debounce != nil {
		filtered := diffs[:0]
		for _, d := range diffs {
			if !p.debounce.IsDuplicate(d) {
				filtered = append(filtered, d)
			}
		}
		diffs = filtered
	}

	// 5. Alert
	if p.alerter != nil && len(diffs) > 0 {
		if err := p.alerter.Notify(diffs); err != nil {
			fmt.Fprintf(p.out, "[warn] alert failed: %v\n", err)
		}
		res.AlertsSent = len(diffs)
		if p.met != nil {
			p.met.RecordAlert(p.host, len(diffs))
		}
	}

	// 6. Persist current snapshot
	if err := p.store.Save(p.host, snap); err != nil {
		return res, fmt.Errorf("pipeline: save state: %w", err)
	}

	// 7. Audit trail
	if p.auditor != nil {
		for _, d := range diffs {
			_ = p.auditor.Record(audit.Entry{
				Host:      p.host,
				Port:      d.Port,
				ChangeType: string(d.Type),
			})
		}
	}

	return res, nil
}
