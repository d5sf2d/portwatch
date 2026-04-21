package probe_test

import (
	"net"
	"testing"
	"time"

	"portwatch/internal/probe"
)

func startTCP(t *testing.T) (port int, stop func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return ln.Addr().(*net.TCPAddr).Port, func() { ln.Close() }
}

func TestProbe_OpenPort(t *testing.T) {
	port, stop := startTCP(t)
	defer stop()
	p := probe.New(time.Second)
	r := p.Probe("127.0.0.1", port)
	if !r.Open {
		t.Fatalf("expected open, got closed: %v", r.Err)
	}
	if r.Latency <= 0 {
		t.Errorf("expected positive latency")
	}
}

func TestProbe_ClosedPort(t *testing.T) {
	p := probe.New(200 * time.Millisecond)
	r := p.Probe("127.0.0.1", 1)
	if r.Open {
		t.Fatal("expected closed")
	}
	if r.Err == nil {
		t.Error("expected non-nil error for closed port")
	}
}

func TestProbeAll_MixedPorts(t *testing.T) {
	port, stop := startTCP(t)
	defer stop()
	p := probe.New(200 * time.Millisecond)
	results := p.ProbeAll("127.0.0.1", []int{port, 1})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	openCount := 0
	for _, r := range results {
		if r.Open {
			openCount++
		}
	}
	if openCount != 1 {
		t.Errorf("expected 1 open, got %d", openCount)
	}
}

func TestProbe_DefaultTimeout(t *testing.T) {
	p := probe.New(0)
	if p == nil {
		t.Fatal("expected non-nil prober")
	}
}
