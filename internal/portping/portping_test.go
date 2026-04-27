package portping_test

import (
	"net"
	"strconv"
	"testing"
	"time"

	"portwatch/internal/portping"
)

func startTCP(t *testing.T) int {
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
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	return port
}

func TestPing_AlivePort(t *testing.T) {
	port := startTCP(t)
	p := portping.New(2*time.Second, 2)
	res := p.Ping("127.0.0.1", port)
	if !res.Alive {
		t.Fatalf("expected alive, got dead: %v", res.Err)
	}
	if res.Latency <= 0 {
		t.Errorf("expected positive latency, got %v", res.Latency)
	}
}

func TestPing_ClosedPort(t *testing.T) {
	p := portping.New(200*time.Millisecond, 1)
	res := p.Ping("127.0.0.1", 1)
	if res.Alive {
		t.Fatal("expected dead port")
	}
	if res.Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestPingAll_MixedPorts(t *testing.T) {
	port := startTCP(t)
	p := portping.New(200*time.Millisecond, 1)
	results := p.PingAll("127.0.0.1", []int{port, 1})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Alive {
		t.Error("first port should be alive")
	}
	if results[1].Alive {
		t.Error("second port should be dead")
	}
}

func TestSummarise_Counts(t *testing.T) {
	results := []portping.Result{
		{Alive: true, Latency: 10 * time.Millisecond},
		{Alive: true, Latency: 20 * time.Millisecond},
		{Alive: false},
	}
	s := portping.Summarise(results)
	if s.Total != 3 {
		t.Errorf("total: want 3, got %d", s.Total)
	}
	if s.Alive != 2 {
		t.Errorf("alive: want 2, got %d", s.Alive)
	}
	if s.Dead != 1 {
		t.Errorf("dead: want 1, got %d", s.Dead)
	}
	if s.AvgRTT != 15*time.Millisecond {
		t.Errorf("avgRTT: want 15ms, got %v", s.AvgRTT)
	}
	if s.MaxRTT != 20*time.Millisecond {
		t.Errorf("maxRTT: want 20ms, got %v", s.MaxRTT)
	}
	if s.MinRTT != 10*time.Millisecond {
		t.Errorf("minRTT: want 10ms, got %v", s.MinRTT)
	}
}

func TestSummarise_Empty(t *testing.T) {
	s := portping.Summarise(nil)
	if s.Total != 0 || s.Alive != 0 {
		t.Errorf("expected zero summary, got %+v", s)
	}
}

func TestAliveResults_Filters(t *testing.T) {
	results := []portping.Result{
		{Alive: true, Port: 80},
		{Alive: false, Port: 81},
		{Alive: true, Port: 443},
	}
	alive := portping.AliveResults(results)
	if len(alive) != 2 {
		t.Fatalf("want 2 alive, got %d", len(alive))
	}
}
