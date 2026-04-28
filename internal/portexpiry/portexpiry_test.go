package portexpiry

import (
	"testing"
	"time"
)

func frozenTracker(at time.Time, defaultMax time.Duration) *Tracker {
	t := New(defaultMax)
	t.now = func() time.Time { return at }
	return t
}

func TestObserve_SetsFirstSeen(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tr := frozenTracker(now, time.Hour)
	tr.Observe("localhost", 80)

	age := tr.Age("localhost", 80)
	if age != 0 {
		t.Fatalf("expected age 0 immediately after observe, got %v", age)
	}
}

func TestObserve_PreservesFirstSeen(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tr := frozenTracker(now, time.Hour)
	tr.Observe("localhost", 80)

	// Advance clock and observe again — first-seen must not change.
	later := now.Add(30 * time.Minute)
	tr.now = func() time.Time { return later }
	tr.Observe("localhost", 80)

	age := tr.Age("localhost", 80)
	if age != 30*time.Minute {
		t.Fatalf("expected 30m age, got %v", age)
	}
}

func TestExpired_WithinLimit(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tr := frozenTracker(now, time.Hour)
	tr.Observe("localhost", 443)

	tr.now = func() time.Time { return now.Add(59 * time.Minute) }
	if tr.Expired("localhost", 443) {
		t.Fatal("expected port not expired within limit")
	}
}

func TestExpired_ExceedsLimit(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tr := frozenTracker(now, time.Hour)
	tr.Observe("localhost", 443)

	tr.now = func() time.Time { return now.Add(61 * time.Minute) }
	if !tr.Expired("localhost", 443) {
		t.Fatal("expected port to be expired after limit")
	}
}

func TestExpired_ZeroMaxAge_NeverExpires(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tr := frozenTracker(now, 0)
	tr.Observe("localhost", 22)

	tr.now = func() time.Time { return now.Add(9999 * time.Hour) }
	if tr.Expired("localhost", 22) {
		t.Fatal("expected port with zero max age to never expire")
	}
}

func TestForget_RemovesPort(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tr := frozenTracker(now, time.Minute)
	tr.Observe("localhost", 8080)
	tr.Forget("localhost", 8080)

	tr.now = func() time.Time { return now.Add(time.Hour) }
	if tr.Expired("localhost", 8080) {
		t.Fatal("forgotten port should not be reported as expired")
	}
	if tr.Age("localhost", 8080) != 0 {
		t.Fatal("forgotten port should have zero age")
	}
}

func TestSetMaxAge_OverridesDefault(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tr := frozenTracker(now, time.Hour)
	tr.Observe("localhost", 3306)
	tr.SetMaxAge("localhost", 3306, 10*time.Minute)

	tr.now = func() time.Time { return now.Add(11 * time.Minute) }
	if !tr.Expired("localhost", 3306) {
		t.Fatal("expected port to expire with overridden shorter max age")
	}
}
