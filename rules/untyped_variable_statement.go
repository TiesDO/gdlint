package rules

import (
	"fmt"

	"github.com/TiesDO/gdlint/core"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

var UntypedVariableStatementRule = core.MatchRule{
	Name:    "untyped_variable_statement",
	Pattern: []byte("(variable_statement !type name: (name) @untyped_var)"),
	Execute: func(match *sitter.QueryMatch, source []byte) ([]core.Warning, error) {
		if len(match.Captures) != 1 {
			return nil, fmt.Errorf("expected only 1 capture, got %d", len(match.Captures))
		}

		capture := match.Captures[0]
		node := capture.Node
		content := string(source[node.StartByte():node.EndByte()])

		message := fmt.Sprintf("variable %s is untyped", content)
		offense := "untyped_variable_statement"

		warnings := make([]core.Warning, 1)
		warnings[0] = *core.NewWarningFromNode(node, message, offense)

		return warnings, nil
	},
}

func init() {
	core.DefaultRuleRegistry.MustRegisterMatchRule(&UntypedVariableStatementRule)
}
