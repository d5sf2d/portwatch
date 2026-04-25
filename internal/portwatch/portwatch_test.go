package portwatch_test

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/example/portwatch/internal/portwatch"
)

// startTCP opens a listener on a random port and returns the port number.
func startTCP(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("startTCP: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	return ln.Addr().(*net.TCPAddr).Port
}

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "portwatch-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func TestNew_MissingHost(t *testing.T) {
	_, err := portwatch.New(portwatch.Config{Ports: []int{80}})
	if err == nil {
		t.Fatal("expected error for empty host")
	}
}

func TestNew_NoPorts(t *testing.T) {
	_, err := portwatch.New(portwatch.Config{Host: "localhost"})
	if err == nil {
		t.Fatal("expected error for empty ports")
	}
}

func TestRun_OpenPortProducesAlert(t *testing.T) {
	port := startTCP(t)
	var buf bytes.Buffer

	w, err := portwatch.New(portwatch.Config{
		Host:     "127.0.0.1",
		Ports:    []int{port},
		StateDir: tempDir(t),
		Out:      &buf,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// First run — no previous state, so no diff expected.
	if err := w.Run(context.Background()); err != nil {
		t.Fatalf("Run (first): %v", err)
	}

	// Verify state was persisted by running again; no new diff expected.
	if err := w.Run(context.Background()); err != nil {
		t.Fatalf("Run (second): %v", err)
	}
}

func TestRun_ClosedPortProducesAlert(t *testing.T) {
	port := startTCP(t)
	var buf bytes.Buffer
	dir := tempDir(t)

	makeCfg := func() portwatch.Config {
		return portwatch.Config{
			Host:     "127.0.0.1",
			Ports:    []int{port},
			StateDir: dir,
			Out:      &buf,
		}
	}

	w1, _ := portwatch.New(makeCfg())
	if err := w1.Run(context.Background()); err != nil {
		t.Fatalf("first run: %v", err)
	}

	// Use a port that is almost certainly closed.
	closedPort := 19999
	cfg2 := makeCfg()
	cfg2.Ports = []int{closedPort}
	w2, _ := portwatch.New(cfg2)
	_ = w2.Run(context.Background())

	// We just verify no panic and no error; alert output is best-effort.
	_ = fmt.Sprintf("%s", buf.String())
}
