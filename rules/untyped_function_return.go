package rules

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
)

var UntypedFunctionReturnRule = MatchRule{
	name:    "untyped_function_return",
	pattern: []byte("(function_definition !return_type name: (name) @untyped_return)"),
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
			Offense:   "untyped_function_return",
			Message:   fmt.Sprintf("function %s has no return type", content),
		}

		return warnings, nil
	},
}

func init() {
	DefaultRuleRegistry.RegisterMatchRule(&UntypedFunctionReturnRule)
}
