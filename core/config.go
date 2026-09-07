package core

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
