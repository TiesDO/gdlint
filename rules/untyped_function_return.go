package rules

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
)

var UntypedFunctionReturnRule = Rule{
	name:    "untyped_function_return",
	pattern: []byte("(function_definition !return_type name: (name) @untyped_return)"),
	execute: func(match *sitter.QueryMatch, source []byte) ([]Warning, error) {
		if len(match.Captures) != 1 {
			return nil, fmt.Errorf("expected only 1 capture, got %d", len(match.Captures))
		}

		capture := match.Captures[0]
		node := capture.Node
		startLine := int(node.StartPoint().Row) + 1
		content := string(source[node.StartByte():node.EndByte()])

		warnings := make([]Warning, 1)
		warnings[0] = Warning{
			LineNumber: startLine,
			Offense:    "untyped_function_return",
			Message:    fmt.Sprintf("function %s has no return type", content),
		}

		return warnings, nil
	},
}

func init() {
	DefaultRuleRegistry.RegisterRule(&UntypedFunctionReturnRule)
}
