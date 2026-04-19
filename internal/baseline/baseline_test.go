package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"portwatch/internal/baseline"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func TestSaveLoad(t *testing.T) {
	m := baseline.New(tempPath(t))
	b := baseline.Capture([]int{80, 443, 8080})
	if err := m.Save(b); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := m.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(got.Ports) != 3 {
		t.Errorf("expected 3 ports, got %d", len(got.Ports))
	}
}

func TestLoad_Missing(t *testing.T) {
	m := baseline.New(filepath.Join(t.TempDir(), "none.json"))
	_, err := m.Load()
	if err != baseline.ErrNoBaseline {
		t.Errorf("expected ErrNoBaseline, got %v", err)
	}
}

func TestCompare_Added(t *testing.T) {
	b := baseline.Capture([]int{80, 443})
	d := baseline.Compare(b, []int{80, 443, 8080})
	if len(d.Added) != 1 || d.Added[0] != 8080 {
		t.Errorf("expected added=[8080], got %v", d.Added)
	}
	if len(d.Removed) != 0 {
		t.Errorf("expected no removed, got %v", d.Removed)
	}
}

func TestCompare_Removed(t *testing.T) {
	b := baseline.Capture([]int{80, 443, 22})
	d := baseline.Compare(b, []int{80, 443})
	if len(d.Removed) != 1 || d.Removed[0] != 22 {
		t.Errorf("expected removed=[22], got %v", d.Removed)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	b := baseline.Capture([]int{80, 443})
	d := baseline.Compare(b, []int{80, 443})
	if d.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestCapture_IndependentCopy(t *testing.T) {
	ports := []int{80, 443}
	b := baseline.Capture(ports)
	ports[0] = 9999
	if b.Ports[0] == 9999 {
		t.Error("Capture should copy the slice")
	}
	os.Unsetenv("_unused")
}
