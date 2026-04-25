package portmute

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/state"
)

func frozen(t time.Time) func() time.Time { return func() time.Time { return t } }

func TestMute_IsMuted_Indefinite(t *testing.T) {
	l := New()
	l.Mute(8080, "maintenance", 0)
	if !l.IsMuted(8080) {
		t.Fatal("expected port 8080 to be muted")
	}
}

func TestMute_Unmute_ClearsEntry(t *testing.T) {
	l := New()
	l.Mute(443, "test", 0)
	l.Unmute(443)
	if l.IsMuted(443) {
		t.Fatal("expected port 443 to be unmuted")
	}
}

func TestMute_Expiry_WithinWindow(t *testing.T) {
	now := time.Now()
	l := New()
	l.now = frozen(now)
	l.Mute(22, "brief", 5*time.Minute)
	l.now = frozen(now.Add(2 * time.Minute))
	if !l.IsMuted(22) {
		t.Fatal("expected port 22 to still be muted within window")
	}
}

func TestMute_Expiry_AfterWindow(t *testing.T) {
	now := time.Now()
	l := New()
	l.now = frozen(now)
	l.Mute(22, "brief", 5*time.Minute)
	l.now = frozen(now.Add(10 * time.Minute))
	if l.IsMuted(22) {
		t.Fatal("expected port 22 to be expired")
	}
}

func TestActive_PrunesExpired(t *testing.T) {
	now := time.Now()
	l := New()
	l.now = frozen(now)
	l.Mute(80, "a", 1*time.Minute)
	l.Mute(443, "b", 0) // indefinite
	l.now = frozen(now.Add(2 * time.Minute))
	actives := l.Active()
	if len(actives) != 1 || actives[0].Port != 443 {
		t.Fatalf("expected only port 443 active, got %v", actives)
	}
}

func TestFilterDiffs_MutedPortRemoved(t *testing.T) {
	l := New()
	l.Mute(8080, "test", 0)
	diffs := []state.Diff{
		{Port: 8080, Status: "opened"},
		{Port: 443, Status: "opened"},
	}
	out := FilterDiffs(diffs, l)
	if len(out) != 1 || out[0].Port != 443 {
		t.Fatalf("expected only port 443 in output, got %v", out)
	}
}

func TestFilterDiffs_NilList(t *testing.T) {
	diffs := []state.Diff{{Port: 22, Status: "opened"}}
	out := FilterDiffs(diffs, nil)
	if len(out) != 1 {
		t.Fatal("nil list should return diffs unchanged")
	}
}

func TestCountMuted(t *testing.T) {
	l := New()
	l.Mute(80, "x", 0)
	l.Mute(443, "y", 0)
	diffs := []state.Diff{
		{Port: 80},
		{Port: 443},
		{Port: 22},
	}
	if n := CountMuted(diffs, l); n != 2 {
		t.Fatalf("expected 2 muted, got %d", n)
	}
}
