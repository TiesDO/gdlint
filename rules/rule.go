package rules

import (
	"fmt"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

type Warning struct {
	StartLine int
	StartChar int
	EndLine   int
	EndChar   int

	Message string
	Offense string
}

func NewWarningFromNode(node sitter.Node, message string, offense string) *Warning {
	return &Warning{
		StartLine: int(node.StartPosition().Row),
		StartChar: int(node.StartPosition().Column),
		EndLine:   int(node.EndPosition().Row),
		EndChar:   int(node.EndPosition().Column),

		Message: message,
		Offense: offense,
	}
}

func (w *Warning) FullMessage() string {
	return fmt.Sprintf("%d:%d (@%s) - %s", w.StartLine+1, w.StartChar+1, w.Offense, w.Message)
}

type Rule interface {
	Name() string
}

type MatchRule struct {
	name    string
	execute func(*sitter.QueryMatch, []byte) ([]Warning, error)
	pattern []byte
}

func (r *MatchRule) Name() string {
	return r.name
}

func (r *MatchRule) Execute(match *sitter.QueryMatch, source []byte) ([]Warning, error) {
	return r.execute(match, source)
}

type NodeRule struct {
	name    string
	execute func(*sitter.Node, []byte) ([]Warning, error)
}

func (r *NodeRule) Name() string {
	return r.name
}

func (r *NodeRule) Execute(node *sitter.Node, source []byte) ([]Warning, error) {
	return r.execute(node, source)
}

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

var DefaultRuleRegistry = NewRuleRegistry()
