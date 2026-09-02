package rules

import (
	"fmt"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

var UntypedConstStatementRule = MatchRule{
	Name:    "untyped_const_statement",
	Pattern: []byte("[(const_statement name: (name) @const_name type: (inferred_type)) (const_statement name: (name) @const_name !type)]"),
	Execute: func(match *sitter.QueryMatch, source []byte) ([]Warning, error) {
		if len(match.Captures) != 1 {
			return nil, fmt.Errorf("expected only 1 capture, got %d", len(match.Captures))
		}

		capture := match.Captures[0]
		node := capture.Node
		content := string(source[node.StartByte():node.EndByte()])

		message := fmt.Sprintf("const %s is untyped", content)
		offense := "untyped_const_statement"

		warnings := make([]Warning, 1)
		warnings[0] = *NewWarningFromNode(node, message, offense)

		return warnings, nil
	},
}

func init() {
	DefaultRuleRegistry.RegisterMatchRule(&UntypedConstStatementRule)
}
