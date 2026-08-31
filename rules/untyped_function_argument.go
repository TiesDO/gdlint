package rules

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
)

var UntypedFunctionArgumentRule = MatchRule{
	name:    "untyped_function_argument",
	pattern: []byte("(function_definition parameters: (parameters (identifier) @untyped_function_argument))"),
	execute: func(match *sitter.QueryMatch, source []byte) ([]Warning, error) {
		if len(match.Captures) != 1 {
			return nil, fmt.Errorf("expected only 1 capture, got %d", len(match.Captures))
		}

		capture := match.Captures[0]
		node := capture.Node
		content := string(source[node.StartByte():node.EndByte()])

		warnings := make([]Warning, 1)
		warnings[0] = Warning{
			StartLine: int(node.StartPoint().Row),
			StartChar: int(node.StartPoint().Column),
			EndLine:   int(node.EndPoint().Row),
			EndChar:   int(node.EndPoint().Column),
			Offense:   "untyped_function_argument",
			Message:   fmt.Sprintf("untyped argument %s", content),
		}

		return warnings, nil
	},
}

func init() {
	DefaultRuleRegistry.RegisterMatchRule(&UntypedFunctionArgumentRule)
}
