package rules

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	gdscript "github.com/prestonknopp/tree-sitter-gdscript/bindings/go"
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

func (r *RuleRunner) RunRules(names []string, ctx context.Context) ([]Warning, error) {
	err := r.UpdateTree(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to parse source into tree: %v", err)
	}

	match_rules := make([]*MatchRule, 0)
	node_rules := make([]*NodeRule, 0)
	for _, name := range names {
		rule := r.registry.GetByName(name)
		if rule == nil {
			return nil, fmt.Errorf("failed to find rule '%s'", name)
		}

		if match_rule, ok := rule.(*MatchRule); ok {
			match_rules = append(match_rules, match_rule)
		}

		if node_rule, ok := rule.(*NodeRule); ok {
			node_rules = append(node_rules, node_rule)
		}
	}

	out_warnings := make([]Warning, 0)

	match_warnings, err := r.runMatchRules(match_rules)

	if err != nil {
		return nil, fmt.Errorf("failed to run match_rules: %v", err)
	}

	if len(match_warnings) > 0 {
		out_warnings = append(out_warnings, match_warnings...)
	}

	node_warnings, err := r.runNodeRules(node_rules)

	if err != nil {
		return nil, fmt.Errorf("failed to run node_rules: %v", err)
	}

	if len(node_warnings) > 0 {
		out_warnings = append(out_warnings, node_warnings...)
	}

	return out_warnings, nil
}

func (r *RuleRunner) UpdateTree(ctx context.Context) error {
	tree := r.parser.ParseCtx(ctx, r.source, r.tree)

	if tree == nil {
		return errors.New("failed to parse tree")
	}

	r.tree = tree
	return nil
}

func (r *RuleRunner) SprintTree() string {
	treeCursor := r.tree.Walk()
	var builder strings.Builder
	r.fprintTreeNode(treeCursor, 0, &builder)

	return builder.String()
}

func (r *RuleRunner) fprintTreeNode(cursor *sitter.TreeCursor, depth int, builder *strings.Builder) {
	indent := strings.Repeat("  ", depth)
	node := cursor.Node()

	prefix := ""
	if fieldName := cursor.FieldName(); fieldName != "" {
		prefix = fieldName + ": "
	}

	start := node.StartPosition()
	end := node.EndPosition()

	fmt.Fprintf(builder,
		"%s%s%s [%d, %d] - [%d, %d]\n",
		indent,
		prefix,
		node.GrammarName(),
		start.Row, start.Column,
		end.Row, end.Column)

	if cursor.GotoFirstChild() {
		r.fprintTreeNode(cursor, depth+1, builder)
		for cursor.GotoNextSibling() {
			r.fprintTreeNode(cursor, depth+1, builder)
		}
		cursor.GotoParent()
	}
}

func (r *RuleRunner) runMatchRules(rules []*MatchRule) ([]Warning, error) {
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

	query, qerr := sitter.NewQuery(r.language, pattern_buf.String())

	if qerr != nil {
		return nil, fmt.Errorf("failed to create query: %v", err)
	}

	cursor := sitter.NewQueryCursor()

	out_warnings := make([]Warning, 0)

	matches := cursor.Matches(query, r.tree.RootNode(), r.source)

	for {
		match := matches.Next()
		if match == nil {
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

func (r *RuleRunner) runNodeRules(rules []*NodeRule) ([]Warning, error) {
	out_warnings := make([]Warning, 0)

	for _, rule := range rules {
		warnings, err := rule.Execute(r.tree.RootNode(), r.source)

		if err != nil {
			return nil, err
		}

		out_warnings = append(out_warnings, warnings...)
	}

	return out_warnings, nil
}

var DefaultRuleRegistry = NewRuleRegistry()
