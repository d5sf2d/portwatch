package policy

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func makeEvaluator(rules []Rule) *Evaluator {
	return New(rules)
}

func TestEvaluate_NoViolation(t *testing.T) {
	e := makeEvaluator([]Rule{
		{Name: "no-telnet", Port: 23, MustBeClosed: true},
	})
	vs := e.Evaluate("localhost", []int{80, 443})
	if len(vs) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(vs))
	}
}

func TestEvaluate_ViolationDetected(t *testing.T) {
	e := makeEvaluator([]Rule{
		{Name: "no-telnet", Port: 23, MustBeClosed: true},
	})
	vs := e.Evaluate("localhost", []int{80, 23})
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
	if vs[0].Port != 23 {
		t.Errorf("expected port 23, got %d", vs[0].Port)
	}
	if vs[0].Rule.Name != "no-telnet" {
		t.Errorf("unexpected rule name %q", vs[0].Rule.Name)
	}
}

func TestEvaluate_HostFilter_Skips(t *testing.T) {
	e := makeEvaluator([]Rule{
		{Name: "db-internal", Port: 5432, MustBeClosed: true, AllowedHosts: []string{"db.internal"}},
	})
	// rule only applies to db.internal; other hosts should be skipped
	vs := e.Evaluate("web.internal", []int{5432})
	if len(vs) != 0 {
		t.Fatalf("expected 0 violations for non-matching host, got %d", len(vs))
	}
}

func TestEvaluate_HostFilter_Matches(t *testing.T) {
	e := makeEvaluator([]Rule{
		{Name: "db-internal", Port: 5432, MustBeClosed: true, AllowedHosts: []string{"db.internal"}},
	})
	vs := e.Evaluate("db.internal", []int{5432})
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
}

func TestViolation_String(t *testing.T) {
	v := Violation{Rule: Rule{Name: "no-ftp"}, Port: 21, Host: "host1"}
	s := v.String()
	if s == "" {
		t.Error("expected non-empty violation string")
	}
}

func TestLoadFile_Valid(t *testing.T) {
	rf := RuleFile{Rules: []Rule{
		{Name: "no-telnet", Port: 23, MustBeClosed: true},
	}}
	b, _ := json.Marshal(rf)
	tmp := filepath.Join(t.TempDir(), "policy.json")
	_ = os.WriteFile(tmp, b, 0o644)

	ev, err := LoadFile(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ev == nil {
		t.Fatal("expected non-nil evaluator")
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	_, err := LoadFile("/nonexistent/policy.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidateRules_MissingName(t *testing.T) {
	err := validateRules([]Rule{{Port: 22, MustBeClosed: true}})
	if err == nil {
		t.Fatal("expected error for empty rule name")
	}
}
