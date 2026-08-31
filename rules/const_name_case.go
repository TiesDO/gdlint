package rules

import (
	"fmt"

	"github.com/TiesDO/gdlint/util"
	sitter "github.com/smacker/go-tree-sitter"
)

var ConstNameCaseRule = Rule{
	name:    "const_name_case",
	pattern: []byte("(const_statement name: (name) @const_name)"),
	execute: func(match *sitter.QueryMatch, source []byte) ([]Warning, error) {
		if len(match.Captures) != 1 {
			return nil, fmt.Errorf("expected only 1 capture, got %d", len(match.Captures))
		}

		capture := match.Captures[0]
		node := capture.Node
		content := string(source[node.StartByte():node.EndByte()])

		if !util.IsConstCase(content) {
			return []Warning{
				Warning{
					StartLine: int(node.StartPoint().Row),
					StartChar: int(node.StartPoint().Column),
					EndLine:   int(node.EndPoint().Row),
					EndChar:   int(node.EndPoint().Column),
					Message:   fmt.Sprintf("const name '%s' must be CONST_CASE", content),
					Offense:   "const_name_case",
				},
			}, nil
		} else {
			return nil, nil
		}
	},
}

func init() {
	DefaultRuleRegistry.RegisterRule(&ConstNameCaseRule)
}
