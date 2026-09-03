package core

import (
	"fmt"
	"log"
	"strings"
)

type DocumentRunner struct {
	registry   *RuleRegistry
	nodeRules  []*NodeRule
	matchRules []*MatchRule
	logger     log.Logger
}

func NewDocumentRunner(registry *RuleRegistry) *DocumentRunner {
	runner := DocumentRunner{
		registry:   registry,
		nodeRules:  []*NodeRule{},
		matchRules: []*MatchRule{},
	}

	return &runner
}

func (r *DocumentRunner) SetRules(rules []string) error {
	r.nodeRules = make([]*NodeRule, 0)
	r.matchRules = make([]*MatchRule, 0)

	for _, ruleName := range rules {
		rule := r.registry.GetByName(ruleName)

		if rule == nil {
			return fmt.Errorf("failed to find rule '%s'", ruleName)
		}

		if nodeRule, ok := rule.(*NodeRule); ok {
			r.nodeRules = append(r.nodeRules, nodeRule)
		} else if matchRule, ok := rule.(*MatchRule); ok {
			r.matchRules = append(r.matchRules, matchRule)
		} else {
			return fmt.Errorf("rule '%s' is of an unsupported type", ruleName)
		}
	}

	return nil
}

func (r *DocumentRunner) CheckDocument(doc *Document) ([]Warning, error) {
	warnings := []Warning{}

	nodeWarnings, err := r.runNodeRules(doc)

	if err != nil {
		return warnings, err
	}

	warnings = append(warnings, nodeWarnings...)

	matchWarnings, err := r.runMatchRules(doc)

	if err != nil {
		return warnings, err
	}

	warnings = append(warnings, matchWarnings...)

	return warnings, nil
}

func (r *DocumentRunner) runNodeRules(doc *Document) ([]Warning, error) {
	warnings := []Warning{}

	for _, rule := range r.nodeRules {
		foundWarnings, err := rule.Execute(doc.parser.tree.RootNode(), doc.source)

		if err != nil {
			return nil, err
		}

		warnings = append(warnings, foundWarnings...)
	}

	return warnings, nil
}

func (r *DocumentRunner) runMatchRules(doc *Document) ([]Warning, error) {
	pattern, err := r.buildMatchRulePattern()

	if err != nil {
		return nil, err
	}

	matches, err := doc.Query(pattern)

	warnings := []Warning{}

	for {
		match := matches.Next()

		if match == nil {
			break
		}

		rule := r.matchRules[match.PatternIndex]

		foundWarnings, err := rule.Execute(match, doc.source)

		if err != nil {
			return nil, err
		}

		warnings = append(warnings, foundWarnings...)
	}

	return warnings, nil
}

func (r *DocumentRunner) buildMatchRulePattern() (string, error) {
	var builder strings.Builder

	for _, rule := range r.matchRules {
		_, err := builder.Write(rule.Pattern)

		if err != nil {
			return "", err
		}

		_, err = builder.WriteString("\n\n")

		if err != nil {
			return "", err
		}
	}

	return builder.String(), nil
}
