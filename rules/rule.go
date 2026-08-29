package rules

import (
	"context"
	"fmt"

	gdscript "github.com/prestonknopp/tree-sitter-gdscript/bindings/go"
	sitter "github.com/smacker/go-tree-sitter"
)

type Warning struct {
	LineNumber int
	Message    string
	Offense    string
}

type RuleOld interface {
	Check(*sitter.Tree, []byte) ([]Warning, error)
}

type Rule struct {
	name    string
	execute func(*sitter.QueryMatch, []byte) ([]Warning, error)
	pattern []byte
}

func (r *Rule) Execute(match *sitter.QueryMatch, source []byte) ([]Warning, error) {
	return r.execute(match, source)
}

type RuleRegistry struct {
	rules map[string]*Rule
}

func NewRuleRegistry() RuleRegistry {
	registry := RuleRegistry{
		rules: make(map[string]*Rule),
	}
	return registry
}

func (r *RuleRegistry) RegisterRule(rule *Rule) error {
	_, exists := r.rules[rule.name]

	if exists {
		return fmt.Errorf("already registered a rule with name '%s'", rule.name)
	}

	r.rules[rule.name] = rule
	return nil
}

func (r *RuleRegistry) GetByName(name string) *Rule {
	rule, ok := r.rules[name]
	if ok {
		return rule
	} else {
		return nil
	}
}

type RuleRunner struct {
	registry *RuleRegistry
	source   []byte
	parser   *sitter.Parser
	language *sitter.Language
	tree     *sitter.Tree
}

func NewRuleRunner(registry *RuleRegistry, source []byte) RuleRunner {
	runner := RuleRunner{
		registry: registry,
		source:   source,
		parser:   sitter.NewParser(),
		language: nil,
		tree:     nil,
	}

	runner.language = sitter.NewLanguage(gdscript.Language())
	runner.parser.SetLanguage(runner.language)

	return runner
}

func (r *RuleRunner) RunRule(name string, ctx context.Context) ([]Warning, error) {
	tree, err := r.parser.ParseCtx(ctx, r.tree, r.source)
	r.tree = tree

	if err != nil {
		return nil, fmt.Errorf("failed to parse source into tree: %v", err)
	}

	rule := r.registry.GetByName(name)

	if rule == nil {
		return nil, fmt.Errorf("failed to find rule '%s' in registry", name)
	}

	query, err := sitter.NewQuery(rule.pattern, r.language)

	if err != nil {
		return nil, fmt.Errorf("failed to instantiate query: %v", err)
	}

	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, r.tree.RootNode())

	out_warnings := make([]Warning, 0)

	for {
		match, ok := cursor.NextMatch()

		if !ok {
			break
		}

		warnings, err := rule.Execute(match, r.source)

		if err != nil {
			return nil, fmt.Errorf("failed to execute rule on match: %v", err)
		}

		out_warnings = append(out_warnings, warnings...)
	}

	return out_warnings, nil
}

// TODO: make this combine the patterns, keep track of match indexes and pass them to the correct rule executor
// func (r *RuleRunner) RunRules(names []string, ctx context.Context) ([]Warning, error) {
// }

var DefaultRuleRegistry = NewRuleRegistry()
