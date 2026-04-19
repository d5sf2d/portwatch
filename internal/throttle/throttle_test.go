package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

func TestAllow_FirstCall(t *testing.T) {
	th := throttle.New(1 * time.Minute)
	if !th.Allow(8080) {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_WithinCooldown(t *testing.T) {
	th := throttle.New(1 * time.Minute)
	th.Allow(8080)
	if th.Allow(8080) {
		t.Fatal("expected second call within cooldown to be denied")
	}
}

func TestAllow_AfterCooldown(t *testing.T) {
	th := throttle.New(10 * time.Millisecond)
	th.Allow(9090)
	time.Sleep(20 * time.Millisecond)
	if !th.Allow(9090) {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestAllow_DifferentPorts(t *testing.T) {
	th := throttle.New(1 * time.Minute)
	th.Allow(80)
	if !th.Allow(443) {
		t.Fatal("expected different port to be allowed independently")
	}
}

func TestReset_ClearsPort(t *testing.T) {
	th := throttle.New(1 * time.Minute)
	th.Allow(3000)
	th.Reset(3000)
	if !th.Allow(3000) {
		t.Fatal("expected allow after reset")
	}
}

func TestResetAll(t *testing.T) {
	th := throttle.New(1 * time.Minute)
	th.Allow(80)
	th.Allow(443)
	th.ResetAll()
	if !th.Allow(80) || !th.Allow(443) {
		t.Fatal("expected both ports allowed after ResetAll")
	}
}

func TestActivePorts_WithinWindow(t *testing.T) {
	th := throttle.New(1 * time.Minute)
	th.Allow(8080)
	th.Allow(9090)
	ports := th.ActivePorts()
	if len(ports) != 2 {
		t.Fatalf("expected 2 active ports, got %d", len(ports))
	}
}

func TestActivePorts_AfterExpiry(t *testing.T) {
	th := throttle.New(10 * time.Millisecond)
	th.Allow(7070)
	time.Sleep(20 * time.Millisecond)
	if len(th.ActivePorts()) != 0 {
		t.Fatal("expected no active ports after cooldown expired")
	}
}
