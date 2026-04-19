package fingerprint

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func startBannerServer(t *testing.T, banner string) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		fmt.Fprint(conn, banner)
		conn.Close()
	}()
	return ln.Addr().(*net.TCPAddr).Port
}

func TestProbe_BannerSSH(t *testing.T) {
	port := startBannerServer(t, "SSH-2.0-OpenSSH_8.9")
	f := New(time.Second)
	svc := f.Probe("127.0.0.1", port)
	if svc.Guess != "ssh" {
		t.Errorf("expected ssh, got %q", svc.Guess)
	}
	if svc.Banner == "" {
		t.Error("expected non-empty banner")
	}
}

func TestProbe_BannerHTTP(t *testing.T) {
	port := startBannerServer(t, "HTTP/1.1 200 OK")
	f := New(time.Second)
	svc := f.Probe("127.0.0.1", port)
	if svc.Guess != "http" {
		t.Errorf("expected http, got %q", svc.Guess)
	}
}

func TestProbe_ClosedPort(t *testing.T) {
	f := New(200 * time.Millisecond)
	svc := f.Probe("127.0.0.1", 1)
	if svc.Guess != "unknown" {
		t.Errorf("expected unknown, got %q", svc.Guess)
	}
}

func TestProbe_FallbackByPort(t *testing.T) {
	port := startBannerServer(t, "") // no banner
	// We can't control what port is assigned, so test guessService directly.
	if g := guessService(3306, ""); g != "mysql" {
		t.Errorf("expected mysql, got %q", g)
	}
	if g := guessService(5432, ""); g != "postgres" {
		t.Errorf("expected postgres, got %q", g)
	}
	_ = port
}

func TestProbeAll_ReturnsAllPorts(t *testing.T) {
	p1 := startBannerServer(t, "SSH-2.0-OpenSSH")
	p2 := startBannerServer(t, "HTTP/1.0")
	f := New(time.Second)
	results := f.ProbeAll("127.0.0.1", []int{p1, p2})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Port != p1 || results[1].Port != p2 {
		t.Error("ports not preserved in order")
	}
}
