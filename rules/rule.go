package rules

import (
	"bufio"
	"bytes"
	"context"
	"fmt"

	gdscript "github.com/prestonknopp/tree-sitter-gdscript/bindings/go"
	sitter "github.com/smacker/go-tree-sitter"
)

type Warning struct {
	StartLine int
	StartChar int
	EndLine   int
	EndChar   int

	Message string
	Offense string
}

type RuleOld interface {
	Check(*sitter.Tree, []byte) ([]Warning, error)
}

type Rule struct {
	name    string
	execute func(*sitter.QueryMatch, []byte) ([]Warning, error)
	pattern []byte
}

func (r *Rule) Name() string {
	return r.name
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

func (r *RuleRegistry) RuleNames() []string {
	names := make([]string, 0)

	for _, rule := range r.rules {
		names = append(names, rule.name)
	}

	return names
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

	if err != nil {
		return nil, fmt.Errorf("failed to parse source into tree: %v", err)
	}

	r.tree = tree

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

func (r *RuleRunner) RunRules(names []string, ctx context.Context) ([]Warning, error) {
	rules := make([]*Rule, 0)
	for _, name := range names {
		rule := r.registry.GetByName(name)
		if rule == nil {
			return nil, fmt.Errorf("failed to find rule '%s'", name)
		}
		rules = append(rules, rule)
	}

	var pattern_buf bytes.Buffer
	pattern_writer := bufio.NewWriter(&pattern_buf)

	for _, rule := range rules {
		_, err := pattern_writer.Write(rule.pattern)

		if err != nil {
			return nil, fmt.Errorf("failed to write pattern to buffer for rule '%s': %v", rule.name, err)
		}

		_, err = pattern_writer.WriteString("\n\n")

		if err != nil {
			return nil, fmt.Errorf("failed to write pattern separator to buffer for rule '%s': %v", rule.name, err)
		}
	}

	err := pattern_writer.Flush()

	if err != nil {
		return nil, fmt.Errorf("failed to flush aggregate pattern writer: %v", err)
	}

	tree, err := r.parser.ParseCtx(context.Background(), r.tree, r.source)

	if err != nil {
		return nil, fmt.Errorf("failed to parse source into tree: %v", err)
	}

	r.tree = tree

	query, err := sitter.NewQuery(pattern_buf.Bytes(), r.language)

	if err != nil {
		return nil, fmt.Errorf("failed to create query: %v", err)
	}

	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, r.tree.RootNode())

	out_warnings := make([]Warning, 0)

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		rule := rules[match.PatternIndex]
		warnings, err := rule.execute(match, r.source)

		if err != nil {
			return nil, err
		}

		out_warnings = append(out_warnings, warnings...)
	}

	return out_warnings, nil
}

var DefaultRuleRegistry = NewRuleRegistry()
