package snapshot_test

import (
	"strings"
	"testing"
	"time"

	"github.com/example/portwatch/internal/snapshot"
)

func makePorts(nums ...int) []snapshot.Port {
	ports := make([]snapshot.Port, len(nums))
	for i, n := range nums {
		ports[i] = snapshot.Port{Number: n, Protocol: "tcp"}
	}
	return ports
}

func TestNew_SortsPorts(t *testing.T) {
	snap := snapshot.New("localhost", makePorts(443, 80, 22))
	if len(snap.Ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(snap.Ports))
	}
	if snap.Ports[0].Number != 22 || snap.Ports[1].Number != 80 || snap.Ports[2].Number != 443 {
		t.Errorf("ports not sorted: %v", snap.Ports)
	}
}

func TestNew_SetsHost(t *testing.T) {
	snap := snapshot.New("myhost", nil)
	if snap.Host != "myhost" {
		t.Errorf("expected host 'myhost', got %q", snap.Host)
	}
}

func TestNew_StampsTime(t *testing.T) {
	before := time.Now().UTC().Add(-time.Second)
	snap := snapshot.New("h", nil)
	after := time.Now().UTC().Add(time.Second)
	if snap.CapturedAt.Before(before) || snap.CapturedAt.After(after) {
		t.Errorf("CapturedAt %v outside expected range", snap.CapturedAt)
	}
}

func TestPortSet_KeyFormat(t *testing.T) {
	snap := snapshot.New("h", []snapshot.Port{
		{Number: 80, Protocol: "tcp", Service: "http"},
	})
	set := snap.PortSet()
	p, ok := set["tcp:80"]
	if !ok {
		t.Fatal("expected key 'tcp:80' in port set")
	}
	if p.Service != "http" {
		t.Errorf("expected service 'http', got %q", p.Service)
	}
}

func TestPortSet_EmptySnapshot(t *testing.T) {
	snap := snapshot.New("h", nil)
	if len(snap.PortSet()) != 0 {
		t.Error("expected empty port set")
	}
}

func TestSummary_ContainsHost(t *testing.T) {
	snap := snapshot.New("scanhost", makePorts(22, 443))
	summary := snap.Summary()
	if !strings.Contains(summary, "scanhost") {
		t.Errorf("summary missing host: %q", summary)
	}
	if !strings.Contains(summary, "ports=2") {
		t.Errorf("summary missing port count: %q", summary)
	}
}
