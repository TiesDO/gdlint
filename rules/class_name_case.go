package rules

import (
	"fmt"

	"github.com/TiesDO/gdlint/util"
	sitter "github.com/smacker/go-tree-sitter"
)

var ClassNameCaseRule = MatchRule{
	name:    "class_name_case",
	pattern: []byte("[(class_name_statement name: (name) @class_name) (class_definition name: (name) @class_definition)]"),
	execute: func(match *sitter.QueryMatch, source []byte) ([]Warning, error) {
		if len(match.Captures) != 1 {
			return nil, fmt.Errorf("expected only 1 capture, got %d", len(match.Captures))
		}

		capture := match.Captures[0]
		node := capture.Node
		content := string(source[node.StartByte():node.EndByte()])

		if !util.IsPascalCase(content) {
			return []Warning{
				{
					StartLine: int(node.StartPoint().Row),
					StartChar: int(node.StartPoint().Column),
					EndLine:   int(node.EndPoint().Row),
					EndChar:   int(node.EndPoint().Column),
					Message:   fmt.Sprintf("class name '%s' must be PascalCase", content),
					Offense:   "class_name_case",
				},
			}, nil
		} else {
			return nil, nil
		}
	},
}

func init() {
	DefaultRuleRegistry.RegisterMatchRule(&ClassNameCaseRule)
}
