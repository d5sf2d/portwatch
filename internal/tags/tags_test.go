package tags_test

import (
	"testing"

	"github.com/user/portwatch/internal/tags"
)

func makeRegistry(t *testing.T) *tags.Registry {
	t.Helper()
	r := tags.New()
	if err := r.Add(tags.Tag{Port: 80, Label: "http", Description: "HTTP traffic"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := r.Add(tags.Tag{Port: 443, Label: "https"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return r
}

func TestAdd_DuplicatePort(t *testing.T) {
	r := makeRegistry(t)
	err := r.Add(tags.Tag{Port: 80, Label: "duplicate"})
	if err == nil {
		t.Fatal("expected error for duplicate port")
	}
}

func TestAdd_EmptyLabel(t *testing.T) {
	r := tags.New()
	err := r.Add(tags.Tag{Port: 8080, Label: ""})
	if err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestGet_Found(t *testing.T) {
	r := makeRegistry(t)
	tag, ok := r.Get(80)
	if !ok {
		t.Fatal("expected tag for port 80")
	}
	if tag.Label != "http" {
		t.Errorf("expected label 'http', got %q", tag.Label)
	}
}

func TestGet_NotFound(t *testing.T) {
	r := makeRegistry(t)
	_, ok := r.Get(9999)
	if ok {
		t.Fatal("expected no tag for port 9999")
	}
}

func TestLabel_FallsBackToDefault(t *testing.T) {
	r := tags.New()
	label := r.Label(3000)
	if label != "port-3000" {
		t.Errorf("expected 'port-3000', got %q", label)
	}
}

func TestRemove_Existing(t *testing.T) {
	r := makeRegistry(t)
	if err := r.Remove(80); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.Get(80); ok {
		t.Fatal("expected port 80 to be removed")
	}
}

func TestRemove_Missing(t *testing.T) {
	r := tags.New()
	if err := r.Remove(1234); err == nil {
		t.Fatal("expected error removing non-existent tag")
	}
}

func TestAll_ReturnsAllTags(t *testing.T) {
	r := makeRegistry(t)
	all := r.All()
	if len(all) != 2 {
		t.Errorf("expected 2 tags, got %d", len(all))
	}
}
