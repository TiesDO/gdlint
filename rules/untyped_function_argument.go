package rules

import (
	"fmt"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

var UntypedFunctionArgumentRule = MatchRule{
	Name:    "untyped_function_argument",
	Pattern: []byte("(function_definition parameters: (parameters (identifier) @untyped_function_argument))"),
	Execute: func(match *sitter.QueryMatch, source []byte) ([]Warning, error) {
		if len(match.Captures) != 1 {
			return nil, fmt.Errorf("expected only 1 capture, got %d", len(match.Captures))
		}

		capture := match.Captures[0]
		node := capture.Node
		content := string(source[node.StartByte():node.EndByte()])

		message := fmt.Sprintf("untyped argument %s", content)
		offense := "untyped_function_argument"

		warnings := make([]Warning, 1)
		warnings[0] = *NewWarningFromNode(node, message, offense)

		return warnings, nil
	},
}

func init() {
	DefaultRuleRegistry.RegisterMatchRule(&UntypedFunctionArgumentRule)
}
