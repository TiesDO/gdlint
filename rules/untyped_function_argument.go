package rules

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
)

var untypedFunctionArgumentRule = Rule{
	name:    "untyped_function_argument",
	pattern: []byte("(function_definition parameters: (parameters (identifier) @untyped_function_argument))"),
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
			Offense:    "untyped_function_argument",
			Message:    fmt.Sprintf("untyped argument %s", content),
		}

		return warnings, nil
	},
}

func init() {
	DefaultRuleRegistry.RegisterRule(&untypedFunctionArgumentRule)
}
