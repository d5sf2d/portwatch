package grouping_test

import (
	"sort"
	"testing"

	"github.com/example/portwatch/internal/grouping"
)

func makeRegistry(t *testing.T) *grouping.Registry {
	t.Helper()
	r := grouping.New()
	if err := r.Add(grouping.Group{Name: "web", Ports: []int{80, 443, 8080}}); err != nil {
		t.Fatalf("unexpected error adding web group: %v", err)
	}
	if err := r.Add(grouping.Group{Name: "db", Ports: []int{5432, 3306}}); err != nil {
		t.Fatalf("unexpected error adding db group: %v", err)
	}
	return r
}

func TestAdd_DuplicateName(t *testing.T) {
	r := makeRegistry(t)
	err := r.Add(grouping.Group{Name: "web", Ports: []int{9000}})
	if err == nil {
		t.Fatal("expected error for duplicate group name, got nil")
	}
}

func TestAdd_EmptyName(t *testing.T) {
	r := grouping.New()
	err := r.Add(grouping.Group{Name: "", Ports: []int{80}})
	if err == nil {
		t.Fatal("expected error for empty group name, got nil")
	}
}

func TestAdd_NoPorts(t *testing.T) {
	r := grouping.New()
	err := r.Add(grouping.Group{Name: "empty", Ports: nil})
	if err == nil {
		t.Fatal("expected error for group with no ports, got nil")
	}
}

func TestGet_Found(t *testing.T) {
	r := makeRegistry(t)
	g, ok := r.Get("web")
	if !ok {
		t.Fatal("expected to find group 'web'")
	}
	if g.Name != "web" {
		t.Errorf("expected name 'web', got %q", g.Name)
	}
}

func TestGet_NotFound(t *testing.T) {
	r := makeRegistry(t)
	_, ok := r.Get("nonexistent")
	if ok {
		t.Fatal("expected group not found, but got ok=true")
	}
}

func TestGroupsForPort_MultipleGroups(t *testing.T) {
	r := grouping.New()
	_ = r.Add(grouping.Group{Name: "web", Ports: []int{80, 443}})
	_ = r.Add(grouping.Group{Name: "all", Ports: []int{80, 443, 5432}})

	names := r.GroupsForPort(80)
	sort.Strings(names)
	if len(names) != 2 || names[0] != "all" || names[1] != "web" {
		t.Errorf("expected [all web], got %v", names)
	}
}

func TestGroupsForPort_NoMatch(t *testing.T) {
	r := makeRegistry(t)
	names := r.GroupsForPort(9999)
	if len(names) != 0 {
		t.Errorf("expected no groups for port 9999, got %v", names)
	}
}

func TestAll_ReturnsAllGroups(t *testing.T) {
	r := makeRegistry(t)
	all := r.All()
	if len(all) != 2 {
		t.Errorf("expected 2 groups, got %d", len(all))
	}
}
