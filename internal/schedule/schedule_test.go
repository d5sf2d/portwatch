package schedule_test

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/schedule"
	"github.com/user/portwatch/internal/state"
)

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func TestRunner_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/state.json"

	sc := scanner.New("127.0.0.1", 50*time.Millisecond)
	st, err := state.NewStore(path)
	if err != nil {
		t.Fatal(err)
	}

	var buf strings.Builder
	n := alert.New(&buf)

	port := freePort(t)

	// Start a listener mid-run to simulate a new open port
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	r := schedule.New(sc, st, n, []int{port}, 60*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	r.Run(ctx) // returns when ctx expires

	if buf.Len() == 0 {
		t.Error("expected alert output, got none")
	}
}

func TestRunner_NoAlertOnFirstRun(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/state.json"

	sc := scanner.New("127.0.0.1", 50*time.Millisecond)
	st, err := state.NewStore(path)
	if err != nil {
		t.Fatal(err)
	}

	var buf strings.Builder
	n := alert.New(&buf)

	r := schedule.New(sc, st, n, []int{os.Getpid()}, 1*time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	r.Run(ctx)

	if buf.Len() != 0 {
		t.Errorf("expected no alert on first run, got: %s", buf.String())
	}
}
