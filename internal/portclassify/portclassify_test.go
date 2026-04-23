package portclassify_test

import (
	"testing"

	"github.com/example/portwatch/internal/portclassify"
)

func TestClassify_SystemPort(t *testing.T) {
	c := portclassify.New()
	r, err := c.Classify(80)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Tier != portclassify.TierSystem {
		t.Errorf("expected system tier, got %s", r.Tier)
	}
	if !r.Privileged {
		t.Error("expected port 80 to be privileged")
	}
}

func TestClassify_RegisteredPort(t *testing.T) {
	c := portclassify.New()
	r, err := c.Classify(8080)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Tier != portclassify.TierRegistered {
		t.Errorf("expected registered tier, got %s", r.Tier)
	}
	if r.Privileged {
		t.Error("expected port 8080 to be unprivileged")
	}
}

func TestClassify_DynamicPort(t *testing.T) {
	c := portclassify.New()
	r, err := c.Classify(55000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Tier != portclassify.TierDynamic {
		t.Errorf("expected dynamic tier, got %s", r.Tier)
	}
}

func TestClassify_BoundaryPort1023(t *testing.T) {
	c := portclassify.New()
	r, err := c.Classify(1023)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Tier != portclassify.TierSystem {
		t.Errorf("port 1023 should be system tier, got %s", r.Tier)
	}
}

func TestClassify_InvalidPort(t *testing.T) {
	c := portclassify.New()
	for _, p := range []int{-1, 65536, 100000} {
		_, err := c.Classify(p)
		if err == nil {
			t.Errorf("expected error for port %d", p)
		}
	}
}

func TestClassifyAll_SkipsInvalid(t *testing.T) {
	c := portclassify.New()
	results := c.ClassifyAll([]int{22, -1, 443, 70000})
	if len(results) != 3 {
		t.Fatalf("expected 3 results (invalid skipped), got %d", len(results))
	}
}

func TestClassifyAll_Empty(t *testing.T) {
	c := portclassify.New()
	results := c.ClassifyAll([]int{})
	if len(results) != 0 {
		t.Errorf("expected empty result slice, got %d", len(results))
	}
}
