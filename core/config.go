package core

type RunnerConfig interface {
	ShouldCheckFile(fileName string) bool
	ShouldCheckRule(ruleName string) bool

	GetEnabledRules() []Rule
}

type DefaultRunnerConfig struct{}

func (c *DefaultRunnerConfig) ShouldCheckFile(_ string) bool {
	return true
}

func (c *DefaultRunnerConfig) ShouldCheckRule(_ string) bool {
	return true
}

func (c *DefaultRunnerConfig) GetEnabledRules() []Rule {
	return DefaultRuleRegistry.GetRules()
}
