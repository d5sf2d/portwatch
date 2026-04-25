package portpriority_test

import (
	"testing"

	"github.com/example/portwatch/internal/portpriority"
)

func makeRegistry(t *testing.T) *portpriority.Registry {
	t.Helper()
	return portpriority.New()
}

func TestGet_CriticalSSH(t *testing.T) {
	r := makeRegistry(t)
	if got := r.Get(22); got != portpriority.Critical {
		t.Errorf("port 22: want Critical, got %s", got)
	}
}

func TestGet_CriticalDatabase(t *testing.T) {
	r := makeRegistry(t)
	for _, port := range []int{3306, 5432, 1433, 27017, 6379} {
		if got := r.Get(port); got != portpriority.Critical {
			t.Errorf("port %d: want Critical, got %s", port, got)
		}
	}
}

func TestGet_HighHTTP(t *testing.T) {
	r := makeRegistry(t)
	for _, port := range []int{80, 443, 8080, 8443} {
		if got := r.Get(port); got != portpriority.High {
			t.Errorf("port %d: want High, got %s", port, got)
		}
	}
}

func TestGet_MediumSystemPort(t *testing.T) {
	r := makeRegistry(t)
	// Port 1000 is system range but not in critical/high list
	if got := r.Get(1000); got != portpriority.Medium {
		t.Errorf("port 1000: want Medium, got %s", got)
	}
}

func TestGet_LowDynamicPort(t *testing.T) {
	r := makeRegistry(t)
	if got := r.Get(51000); got != portpriority.Low {
		t.Errorf("port 51000: want Low, got %s", got)
	}
}

func TestSet_OverridesDefault(t *testing.T) {
	r := makeRegistry(t)
	if err := r.Set(9999, portpriority.Critical); err != nil {
		t.Fatalf("Set: unexpected error: %v", err)
	}
	if got := r.Get(9999); got != portpriority.Critical {
		t.Errorf("port 9999: want Critical after override, got %s", got)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	r := makeRegistry(t)
	if err := r.Set(0, portpriority.High); err == nil {
		t.Error("expected error for port 0, got nil")
	}
	if err := r.Set(65536, portpriority.High); err == nil {
		t.Error("expected error for port 65536, got nil")
	}
}

func TestLevel_String(t *testing.T) {
	cases := []struct {
		level portpriority.Level
		want  string
	}{
		{portpriority.Low, "low"},
		{portpriority.Medium, "medium"},
		{portpriority.High, "high"},
		{portpriority.Critical, "critical"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("Level(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}
