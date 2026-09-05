package rules

import (
	"fmt"

	"github.com/TiesDO/gdlint/core"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

var UntypedFunctionReturnRule = core.MatchRule{
	Name:    "untyped_function_return",
	Pattern: []byte("(function_definition !return_type name: (name) @untyped_return)"),
	Execute: func(match *sitter.QueryMatch, _ *sitter.Query, document *core.Document) ([]core.Warning, error) {
		if len(match.Captures) != 1 {
			return nil, fmt.Errorf("expected only 1 capture, got %d", len(match.Captures))
		}

		capture := match.Captures[0]
		node := capture.Node
		content := document.ContentForNode(&node)

		message := fmt.Sprintf("function %s has no return type", content)
		offense := "untyped_function_return"

		warnings := make([]core.Warning, 1)
		warnings[0] = *core.NewWarningFromNode(node, message, offense)

		return warnings, nil
	},
}
