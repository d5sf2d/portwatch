package portretry

import (
	"testing"
	"time"
)

func noBackoff(r *Retryer) {
	r.clock = func(time.Duration) {}
}

func TestProbe_OpenOnFirstAttempt(t *testing.T) {
	r := New(3, 0)
	noBackoff(r)

	result := r.Probe(80, func(port int) bool { return true })

	if !result.Open {
		t.Fatal("expected port to be open")
	}
	if result.Attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", result.Attempts)
	}
}

func TestProbe_OpenOnSecondAttempt(t *testing.T) {
	r := New(3, time.Millisecond)
	noBackoff(r)

	calls := 0
	result := r.Probe(443, func(port int) bool {
		calls++
		return calls >= 2
	})

	if !result.Open {
		t.Fatal("expected port to be open")
	}
	if result.Attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", result.Attempts)
	}
}

func TestProbe_NeverOpen(t *testing.T) {
	r := New(3, 0)
	noBackoff(r)

	result := r.Probe(9999, func(port int) bool { return false })

	if result.Open {
		t.Fatal("expected port to be closed")
	}
	if result.Attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", result.Attempts)
	}
}

func TestNew_ClampsMinAttempts(t *testing.T) {
	r := New(0, 0)
	if r.MaxAttempts != 1 {
		t.Fatalf("expected MaxAttempts=1, got %d", r.MaxAttempts)
	}
}

func TestProbeAll_ReturnsAllResults(t *testing.T) {
	r := New(2, 0)
	noBackoff(r)

	openPorts := map[int]bool{22: true, 80: true}
	results := r.ProbeAll([]int{22, 80, 9999}, func(port int) bool {
		return openPorts[port]
	})

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, res := range results {
		switch res.Port {
		case 22, 80:
			if !res.Open {
				t.Errorf("port %d should be open", res.Port)
			}
		case 9999:
			if res.Open {
				t.Errorf("port %d should be closed", res.Port)
			}
		}
	}
}

func TestProbe_BackoffCalledBetweenAttempts(t *testing.T) {
	r := New(3, 50*time.Millisecond)
	sleptCount := 0
	r.clock = func(d time.Duration) { sleptCount++ }

	r.Probe(8080, func(port int) bool { return false })

	// 3 attempts → 2 sleeps between them
	if sleptCount != 2 {
		t.Fatalf("expected 2 backoff sleeps, got %d", sleptCount)
	}
}
