package healthcheck_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/healthcheck"
)

// startServer opens a TCP listener on a random port and returns its port number.
func startServer(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return ln.Addr().(*net.TCPAddr).Port
}

func TestCheck_OpenPort(t *testing.T) {
	port := startServer(t)
	c := healthcheck.New(time.Second)
	st := c.Check("127.0.0.1", port)

	if !st.Healthy {
		t.Fatalf("expected healthy, got reason: %s", st.Reason)
	}
	if st.Port != port {
		t.Errorf("port mismatch: got %d want %d", st.Port, port)
	}
	if st.LatencyMs < 0 {
		t.Errorf("negative latency: %d", st.LatencyMs)
	}
}

func TestCheck_ClosedPort(t *testing.T) {
	c := healthcheck.New(200 * time.Millisecond)
	st := c.Check("127.0.0.1", 1) // port 1 should be refused

	if st.Healthy {
		t.Fatal("expected unhealthy for closed port")
	}
	if st.Reason == "" {
		t.Error("expected non-empty reason")
	}
}

func TestCheckAll_MixedPorts(t *testing.T) {
	port := startServer(t)
	c := healthcheck.New(500 * time.Millisecond)
	results := c.CheckAll("127.0.0.1", []int{port, 1})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Healthy {
		t.Errorf("port %d should be healthy", port)
	}
	if results[1].Healthy {
		t.Error("port 1 should be unhealthy")
	}
}

func TestCheck_DefaultTimeout(t *testing.T) {
	c := healthcheck.New(0) // should default to 2s
	_ = fmt.Sprintf("%+v", c)  // ensure non-nil
	port := startServer(t)
	st := c.Check("127.0.0.1", port)
	if !st.Healthy {
		t.Fatalf("expected healthy: %s", st.Reason)
	}
}
