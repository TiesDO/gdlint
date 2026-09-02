package rules

import (
	"fmt"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

var UntypedVariableStatementRule = MatchRule{
	name:    "untyped_variable_statement",
	pattern: []byte("(variable_statement !type name: (name) @untyped_var)"),
	execute: func(match *sitter.QueryMatch, source []byte) ([]Warning, error) {
		if len(match.Captures) != 1 {
			return nil, fmt.Errorf("expected only 1 capture, got %d", len(match.Captures))
		}

		capture := match.Captures[0]
		node := capture.Node
		content := string(source[node.StartByte():node.EndByte()])

		message := fmt.Sprintf("variable %s is untyped", content)
		offense := "untyped_variable_statement"

		warnings := make([]Warning, 1)
		warnings[0] = *NewWarningFromNode(node, message, offense)

		return warnings, nil
	},
}

func init() {
	DefaultRuleRegistry.RegisterMatchRule(&UntypedVariableStatementRule)
}
