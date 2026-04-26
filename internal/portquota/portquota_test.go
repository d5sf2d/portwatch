package portquota_test

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/portquota"
	"github.com/example/portwatch/internal/state"
)

func makeSnapshot(host string, ports ...int) state.Snapshot {
	ps := make([]state.PortEntry, len(ports))
	for i, p := range ports {
		ps[i] = state.PortEntry{Port: p, Proto: "tcp"}
	}
	return state.Snapshot{Host: host, Ports: ps, At: time.Now()}
}

func TestCheck_WithinLimit(t *testing.T) {
	q := portquota.New(5)
	if err := q.Check("host1", 5); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCheck_ExceedsLimit(t *testing.T) {
	q := portquota.New(3)
	if err := q.Check("host1", 4); err == nil {
		t.Fatal("expected error when open count exceeds limit")
	}
}

func TestCheck_UnlimitedDefault(t *testing.T) {
	q := portquota.New(0)
	if err := q.Check("host1", 9999); err != nil {
		t.Fatalf("unlimited quota should never error, got %v", err)
	}
}

func TestSetHost_OverridesDefault(t *testing.T) {
	q := portquota.New(10)
	q.SetHost("special", 2)
	if err := q.Check("special", 3); err == nil {
		t.Fatal("expected error for host-specific limit")
	}
	// Other hosts still use default.
	if err := q.Check("other", 9); err != nil {
		t.Fatalf("default limit should allow 9 ports, got %v", err)
	}
}

func TestSetHost_ZeroRemovesOverride(t *testing.T) {
	q := portquota.New(2)
	q.SetHost("h", 100)
	q.SetHost("h", 0) // remove override
	if q.Limit("h") != 2 {
		t.Fatalf("expected default limit 2 after removing override, got %d", q.Limit("h"))
	}
}

func TestExceeded_ReportsCorrectly(t *testing.T) {
	q := portquota.New(3)
	if q.Exceeded("h", 3) {
		t.Fatal("3 ports with limit 3 should not be exceeded")
	}
	if !q.Exceeded("h", 4) {
		t.Fatal("4 ports with limit 3 should be exceeded")
	}
}

func TestExceedingPorts_ReturnsExcessPorts(t *testing.T) {
	q := portquota.New(2)
	snap := makeSnapshot("host", 22, 80, 443, 8080)
	excess := portquota.ExceedingPorts(q, snap)
	if len(excess) != 2 {
		t.Fatalf("expected 2 excess ports, got %d", len(excess))
	}
	if excess[0] != 443 || excess[1] != 8080 {
		t.Fatalf("unexpected excess ports: %v", excess)
	}
}

func TestExceedingPorts_NilQuota(t *testing.T) {
	snap := makeSnapshot("host", 22, 80)
	if portquota.ExceedingPorts(nil, snap) != nil {
		t.Fatal("nil quota should return nil")
	}
}

func TestCountExceeding_BelowLimit(t *testing.T) {
	q := portquota.New(10)
	snap := makeSnapshot("host", 22, 80)
	if n := portquota.CountExceeding(q, snap); n != 0 {
		t.Fatalf("expected 0 excess, got %d", n)
	}
}
