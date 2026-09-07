package core

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gobwas/glob"
)

type RunnerConfig interface {
	ShouldCheckFile(fileName string) bool
	ShouldCheckRule(ruleName string) bool

	GetEnabledRules() []Rule
}

type AllRuleRunnerConfig struct {
	rules []Rule
}

func NewAllRuleRunnerConfig() *AllRuleRunnerConfig {
	return &AllRuleRunnerConfig{rules: DefaultRuleRegistry.GetRules()}
}

func (c *AllRuleRunnerConfig) ShouldCheckFile(fileName string) bool {
	return strings.HasSuffix(fileName, ".gd")
}

func (c *AllRuleRunnerConfig) ShouldCheckRule(ruleName string) bool {
	return slices.ContainsFunc(c.rules, func(rule Rule) bool {
		return rule.Identifier() == ruleName
	})
}

func (c *AllRuleRunnerConfig) GetEnabledRules() []Rule {
	return c.rules
}

type RawConfig struct {
	EnabledRules  []string `json:"enabledRules"`
	DisabledRules []string `json:"disabledRules"`

	IncludedFiles []string `json:"includedFiles"`
	ExcludedFiles []string `json:"excludedFiles"`
}

type StandardConfig struct {
	enabledRules         []Rule
	includedFilePatterns []*glob.Pattern
	excludedFilePatterns []*glob.Pattern
}

func NewStandardConfigFromRaw(raw *RawConfig) (*StandardConfig, error) {
	config := StandardConfig{}

	if len(raw.EnabledRules) == 0 && len(raw.DisabledRules) == 0 {
		// no preferences specified; enable all rules
		config.enabledRules = DefaultRuleRegistry.GetRules()
	} else {
		enabledRuleNames := raw.EnabledRules

		if len(enabledRuleNames) == 0 {
			enabledRuleNames = DefaultRuleRegistry.RuleNames()
		}

		for _, ruleName := range enabledRuleNames {
			if slices.Contains(raw.DisabledRules, ruleName) {
				continue
			}

			rule := DefaultRuleRegistry.GetByName(ruleName)

			if rule == nil {
				return nil, fmt.Errorf("rule %q does not exist", ruleName)
			}

			config.enabledRules = append(config.enabledRules, rule)
		}
	}

	for _, pattern := range raw.IncludedFiles {
		targetPattern := pattern

		// support recursive patterns at the root as well
		if cut, ok := strings.CutPrefix(pattern, "**/"); ok {
			targetPattern = fmt.Sprintf("{%s,%s}", cut, pattern)
		}

		compiled, err := glob.Compile(targetPattern, '/')

		if err != nil {
			return nil, fmt.Errorf("failed to compile include pattern %q: %w", compiled, err)
		}

		config.includedFilePatterns = append(config.includedFilePatterns, compiled)
	}

	for _, pattern := range raw.ExcludedFiles {
		targetPattern := pattern

		// support recursive patterns at the root as well
		if cut, ok := strings.CutPrefix(pattern, "**/"); ok {
			targetPattern = fmt.Sprintf("{%s,%s}", cut, pattern)
		}

		compiled, err := glob.Compile(targetPattern, '/')

		if err != nil {
			return nil, fmt.Errorf("failed to compile exclude pattern %q: %w", compiled, err)
		}

		config.excludedFilePatterns = append(config.excludedFilePatterns, compiled)
	}

	return &config, nil
}

func (c *StandardConfig) ShouldCheckFile(fileName string) bool {
	for _, pattern := range c.includedFilePatterns {
		if !pattern.Match(fileName) {
			return false
		}
	}

	for _, pattern := range c.excludedFilePatterns {
		if pattern.Match(fileName) {
			return false
		}
	}

	return true
}

func (c *StandardConfig) ShouldCheckRule(ruleName string) bool {
	return slices.ContainsFunc(c.enabledRules, func(rule Rule) bool {
		return rule.Identifier() == ruleName
	})
}

func (c *StandardConfig) GetEnabledRules() []Rule {
	return c.enabledRules
}
