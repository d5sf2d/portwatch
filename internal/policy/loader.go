package policy

import (
	"encoding/json"
	"fmt"
	"os"
)

// RuleFile is the JSON schema for a policy file.
type RuleFile struct {
	Rules []Rule `json:"rules"`
}

// LoadFile reads a JSON policy file from path and returns an Evaluator.
func LoadFile(path string) (*Evaluator, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("policy: open %s: %w", path, err)
	}
	defer f.Close()

	var rf RuleFile
	if err := json.NewDecoder(f).Decode(&rf); err != nil {
		return nil, fmt.Errorf("policy: decode %s: %w", path, err)
	}

	if err := validateRules(rf.Rules); err != nil {
		return nil, err
	}

	return New(rf.Rules), nil
}

func validateRules(rules []Rule) error {
	for i, r := range rules {
		if r.Name == "" {
			return fmt.Errorf("policy: rule[%d] missing name", i)
		}
		if r.Port < 0 || r.Port > 65535 {
			return fmt.Errorf("policy: rule %q has invalid port %d", r.Name, r.Port)
		}
	}
	return nil
}
