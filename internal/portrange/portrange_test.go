package portrange

import (
	"testing"
)

func TestParse_SinglePort(t *testing.T) {
	ports, err := Parse("80")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 1 || ports[0] != 80 {
		t.Fatalf("expected [80], got %v", ports)
	}
}

func TestParse_CommaSeparated(t *testing.T) {
	ports, err := Parse("22,80,443")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int{22, 80, 443}
	if len(ports) != len(want) {
		t.Fatalf("expected %v, got %v", want, ports)
	}
	for i, p := range want {
		if ports[i] != p {
			t.Errorf("index %d: expected %d, got %d", i, p, ports[i])
		}
	}
}

func TestParse_Range(t *testing.T) {
	ports, err := Parse("8080-8083")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int{8080, 8081, 8082, 8083}
	if len(ports) != len(want) {
		t.Fatalf("expected %v, got %v", want, ports)
	}
	for i, p := range want {
		if ports[i] != p {
			t.Errorf("index %d: expected %d, got %d", i, p, ports[i])
		}
	}
}

func TestParse_MixedExpr(t *testing.T) {
	ports, err := Parse("22, 8080-8082, 443")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int{22, 443, 8080, 8081, 8082}
	if len(ports) != len(want) {
		t.Fatalf("expected %v, got %v", want, ports)
	}
}

func TestParse_Deduplication(t *testing.T) {
	ports, err := Parse("80,80,80")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 1 || ports[0] != 80 {
		t.Fatalf("expected [80], got %v", ports)
	}
}

func TestParse_EmptyExprReturnsError(t *testing.T) {
	_, err := Parse("")
	if err == nil {
		t.Fatal("expected error for empty expression")
	}
}

func TestParse_InvalidPortReturnsError(t *testing.T) {
	_, err := Parse("abc")
	if err == nil {
		t.Fatal("expected error for non-numeric token")
	}
}

func TestParse_OutOfRangeReturnsError(t *testing.T) {
	_, err := Parse("99999")
	if err == nil {
		t.Fatal("expected error for out-of-range port")
	}
}

func TestParse_ReversedRangeReturnsError(t *testing.T) {
	_, err := Parse("9000-8000")
	if err == nil {
		t.Fatal("expected error for reversed range")
	}
}

func TestContains_Found(t *testing.T) {
	ok, err := Contains("80,443,8080-8090", 8085)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected port 8085 to be contained")
	}
}

func TestContains_NotFound(t *testing.T) {
	ok, err := Contains("80,443", 8080)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected port 8080 to not be contained")
	}
}
