package rules

import "fmt"

var DefaultRuleRegistry = NewRuleRegistry()

type RuleRegistry struct {
	match_rules map[string]*MatchRule
	node_rules  map[string]*NodeRule
}

func NewRuleRegistry() RuleRegistry {
	registry := RuleRegistry{
		match_rules: make(map[string]*MatchRule),
		node_rules:  make(map[string]*NodeRule),
	}
	return registry
}

func (r *RuleRegistry) RuleExists(name string) bool {
	if _, exists := r.match_rules[name]; exists {
		return true
	}

	if _, exists := r.node_rules[name]; exists {
		return true
	}

	return false
}

func (r *RuleRegistry) RegisterMatchRule(rule *MatchRule) error {
	exists := r.RuleExists(rule.name)

	if exists {
		return fmt.Errorf("already registered a rule with name '%s'", rule.name)
	}

	r.match_rules[rule.name] = rule
	return nil
}

func (r *RuleRegistry) RegisterNodeRule(rule *NodeRule) error {
	exists := r.RuleExists(rule.name)

	if exists {
		return fmt.Errorf("already registered a rule with name '%s'", rule.name)
	}

	r.node_rules[rule.name] = rule
	return nil
}

func (r *RuleRegistry) GetByName(name string) Rule {
	if rule, ok := r.match_rules[name]; ok {
		return rule
	}

	if rule, ok := r.node_rules[name]; ok {
		return rule
	}

	return nil
}

func (r *RuleRegistry) RuleNames() []string {
	names := make([]string, 0)

	for _, rule := range r.match_rules {
		names = append(names, rule.name)
	}

	for _, rule := range r.node_rules {
		names = append(names, rule.name)
	}

	return names
}
