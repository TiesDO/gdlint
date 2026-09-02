package rules

import (
	"fmt"

	"github.com/TiesDO/gdlint/util"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

var ClassNameCaseRule = MatchRule{
	Name:    "class_name_case",
	Pattern: []byte("[(class_name_statement name: (name) @class_name) (class_definition name: (name) @class_definition)]"),
	Execute: func(match *sitter.QueryMatch, source []byte) ([]Warning, error) {
		if len(match.Captures) != 1 {
			return nil, fmt.Errorf("expected only 1 capture, got %d", len(match.Captures))
		}

		capture := match.Captures[0]
		node := capture.Node
		content := string(source[node.StartByte():node.EndByte()])

		message := fmt.Sprintf("class name '%s' must be PascalCase", content)
		offense := "class_name_case"

		if !util.IsPascalCase(content) {
			return []Warning{
				*NewWarningFromNode(node, message, offense),
			}, nil
		} else {
			return nil, nil
		}
	},
}

func init() {
	DefaultRuleRegistry.RegisterMatchRule(&ClassNameCaseRule)
}
