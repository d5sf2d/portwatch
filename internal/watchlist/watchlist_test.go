package watchlist_test

import (
	"testing"

	"github.com/user/portwatch/internal/watchlist"
)

func makeWatchlist() *watchlist.Watchlist {
	return &watchlist.Watchlist{
		Groups: []watchlist.Group{
			{Name: "web", Ports: []int{80, 443, 8080}},
			{Name: "database", Ports: []int{5432, 3306}},
		},
	}
}

func TestAllPorts_Deduplication(t *testing.T) {
	wl := &watchlist.Watchlist{
		Groups: []watchlist.Group{
			{Name: "a", Ports: []int{80, 443}},
			{Name: "b", Ports: []int{443, 8080}},
		},
	}
	ports := wl.AllPorts()
	if len(ports) != 3 {
		t.Fatalf("expected 3 unique ports, got %d", len(ports))
	}
	if ports[0] != 80 || ports[1] != 443 || ports[2] != 8080 {
		t.Errorf("unexpected port order: %v", ports)
	}
}

func TestValidate_Valid(t *testing.T) {
	if err := makeWatchlist().Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_EmptyName(t *testing.T) {
	wl := &watchlist.Watchlist{
		Groups: []watchlist.Group{{Name: "", Ports: []int{80}}},
	}
	if err := wl.Validate(); err == nil {
		t.Fatal("expected error for empty group name")
	}
}

func TestValidate_DuplicateName(t *testing.T) {
	wl := &watchlist.Watchlist{
		Groups: []watchlist.Group{
			{Name: "web", Ports: []int{80}},
			{Name: "web", Ports: []int{443}},
		},
	}
	if err := wl.Validate(); err == nil {
		t.Fatal("expected error for duplicate group name")
	}
}

func TestValidate_InvalidPort(t *testing.T) {
	wl := &watchlist.Watchlist{
		Groups: []watchlist.Group{{Name: "bad", Ports: []int{0}}},
	}
	if err := wl.Validate(); err == nil {
		t.Fatal("expected error for out-of-range port")
	}
}

func TestFindGroups(t *testing.T) {
	wl := makeWatchlist()
	groups := wl.FindGroups(443)
	if len(groups) != 1 || groups[0] != "web" {
		t.Errorf("expected [web], got %v", groups)
	}
	if len(wl.FindGroups(9999)) != 0 {
		t.Error("expected no groups for unknown port")
	}
}
