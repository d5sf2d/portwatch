package resolve

import (
	"testing"
)

func TestName_Builtin(t *testing.T) {
	r := New(nil)
	if got := r.Name(22); got != "ssh" {
		t.Fatalf("expected ssh, got %s", got)
	}
}

func TestName_Override(t *testing.T) {
	r := New(map[int]string{9000: "my-service"})
	if got := r.Name(9000); got != "my-service" {
		t.Fatalf("expected my-service, got %s", got)
	}
}

func TestName_OverrideTakesPrecedence(t *testing.T) {
	r := New(map[int]string{22: "custom-ssh"})
	if got := r.Name(22); got != "custom-ssh" {
		t.Fatalf("expected custom-ssh, got %s", got)
	}
}

func TestName_Unknown(t *testing.T) {
	r := New(nil)
	if got := r.Name(19999); got != "port/19999" {
		t.Fatalf("expected port/19999, got %s", got)
	}
}

func TestKnown_True(t *testing.T) {
	r := New(nil)
	if !r.Known(80) {
		t.Fatal("expected port 80 to be known")
	}
}

func TestKnown_False(t *testing.T) {
	r := New(nil)
	if r.Known(19999) {
		t.Fatal("expected port 19999 to be unknown")
	}
}

func TestKnown_ViaOverride(t *testing.T) {
	r := New(map[int]string{9999: "custom"})
	if !r.Known(9999) {
		t.Fatal("expected override port to be known")
	}
}
