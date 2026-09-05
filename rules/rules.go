package rules

import "github.com/TiesDO/gdlint/core"

func RegisterAll(r *core.RuleRegistry) {
	r.MustRegisterMatchRule(&ClassNameCaseRule)
	r.MustRegisterMatchRule(&ConstNameCaseRule)
	r.MustRegisterMatchRule(&UntypedConstStatementRule)
	r.MustRegisterMatchRule(&UntypedFunctionArgumentRule)
	r.MustRegisterMatchRule(&UntypedFunctionReturnRule)
	r.MustRegisterMatchRule(&UntypedVariableStatementRule)

	r.MustRegisterNodeRule(&CodeOrderRule)
}
