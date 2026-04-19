package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/schedule"
	"github.com/user/portwatch/internal/state"
)

func main() {
	host := flag.String("host", "127.0.0.1", "host to scan")
	portsFlag := flag.String("ports", "80,443,8080", "comma-separated list of ports")
	interval := flag.Duration("interval", 30*time.Second, "scan interval")
	timeout := flag.Duration("timeout", 500*time.Millisecond, "per-port dial timeout")
	stateFile := flag.String("state", "/tmp/portwatch_state.json", "path to state file")
	flag.Parse()

	ports, err := parsePorts(*portsFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid ports: %v\n", err)
		os.Exit(1)
	}

	sc := scanner.New(*host, *timeout)

	st, err := state.NewStore(*stateFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "state store error: %v\n", err)
		os.Exit(1)
	}

	n := alert.New(os.Stdout)
	r := schedule.New(sc, st, n, ports, *interval)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Printf("portwatch started — host=%s ports=%v interval=%s\n", *host, ports, *interval)
	if err := r.Run(ctx); err != nil && err != context.Canceled {
		fmt.Fprintf(os.Stderr, "runner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("portwatch stopped")
}

func parsePorts(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	ports := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		v, err := strconv.Atoi(p)
		if err != nil || v < 1 || v > 65535 {
			return nil, fmt.Errorf("invalid port %q", p)
		}
		ports = append(ports, v)
	}
	return ports, nil
}
