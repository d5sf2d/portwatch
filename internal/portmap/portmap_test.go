package portmap_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/portmap"
)

func TestLookup_WellKnownPort(t *testing.T) {
	r := portmap.New()
	name, ok := r.Lookup(22)
	if !ok {
		t.Fatal("expected port 22 to be found")
	}
	if name != "ssh" {
		t.Fatalf("expected ssh, got %s", name)
	}
}

func TestLookup_UnknownPort(t *testing.T) {
	r := portmap.New()
	_, ok := r.Lookup(9999)
	if ok {
		t.Fatal("expected port 9999 to be absent")
	}
}

func TestAdd_OverridesExisting(t *testing.T) {
	r := portmap.New()
	r.Add(80, "my-http")
	name, ok := r.Lookup(80)
	if !ok {
		t.Fatal("expected port 80 to be found after override")
	}
	if name != "my-http" {
		t.Fatalf("expected my-http, got %s", name)
	}
}

func TestAdd_CustomPort(t *testing.T) {
	r := portmap.New()
	r.Add(9200, "elasticsearch")
	name, ok := r.Lookup(9200)
	if !ok {
		t.Fatal("expected custom port to be found")
	}
	if name != "elasticsearch" {
		t.Fatalf("expected elasticsearch, got %s", name)
	}
}

func TestAdd_EmptyNameIgnored(t *testing.T) {
	r := portmap.New()
	r.Add(9200, "")
	_, ok := r.Lookup(9200)
	if ok {
		t.Fatal("expected empty-name add to be ignored")
	}
}

func TestLookupDefault_Found(t *testing.T) {
	r := portmap.New()
	got := r.LookupDefault(443, "unknown")
	if got != "https" {
		t.Fatalf("expected https, got %s", got)
	}
}

func TestLookupDefault_Fallback(t *testing.T) {
	r := portmap.New()
	got := r.LookupDefault(9999, "unknown")
	if got != "unknown" {
		t.Fatalf("expected fallback 'unknown', got %s", got)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	r := portmap.New()
	all := r.All()
	if len(all) == 0 {
		t.Fatal("expected non-empty map from All")
	}
	// Mutating the copy must not affect the registry.
	all[22] = "tampered"
	name, _ := r.Lookup(22)
	if name != "ssh" {
		t.Fatal("All() returned a live reference, not a copy")
	}
}
