package scanner_test

import (
	"net"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

// startTCPServer opens a local TCP listener and returns its port and a close func.
func startTCPServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestScan_OpenPort(t *testing.T) {
	port, close := startTCPServer(t)
	defer close()

	s := scanner.New()
	result, err := s.Scan("127.0.0.1", []int{port})
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	if len(result.Ports) != 1 {
		t.Fatalf("expected 1 port result, got %d", len(result.Ports))
	}
	if !result.Ports[0].Open {
		t.Errorf("expected port %d to be open", port)
	}
}

func TestScan_ClosedPort(t *testing.T) {
	// Port 1 is almost certainly closed in test environments.
	s := scanner.New()
	result, err := s.Scan("127.0.0.1", []int{1})
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	if len(result.Ports) != 1 {
		t.Fatalf("expected 1 port result, got %d", len(result.Ports))
	}
	if result.Ports[0].Open {
		t.Errorf("expected port 1 to be closed")
	}
}

func TestScan_ResultMetadata(t *testing.T) {
	s := scanner.New()
	result, err := s.Scan("127.0.0.1", []int{1, 2})
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if result.Host != "127.0.0.1" {
		t.Errorf("unexpected host: %s", result.Host)
	}
	if result.ScannedAt.IsZero() {
		t.Error("ScannedAt should not be zero")
	}
	if len(result.Ports) != 2 {
		t.Errorf("expected 2 port results, got %d", len(result.Ports))
	}
}
