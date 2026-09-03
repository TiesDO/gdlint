package rules

import (
	"fmt"

	"github.com/TiesDO/gdlint/core"
	"github.com/TiesDO/gdlint/util"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

var ConstNameCaseRule = core.MatchRule{
	Name:    "const_name_case",
	Pattern: []byte("(const_statement name: (name) @const_name)"),
	Execute: func(match *sitter.QueryMatch, source []byte) ([]core.Warning, error) {
		if len(match.Captures) != 1 {
			return nil, fmt.Errorf("expected only 1 capture, got %d", len(match.Captures))
		}

		capture := match.Captures[0]
		node := capture.Node
		content := string(source[node.StartByte():node.EndByte()])

		message := fmt.Sprintf("const name '%s' must be CONST_CASE", content)
		offense := "const_name_case"

		if !util.IsConstCase(content) {
			return []core.Warning{
				*core.NewWarningFromNode(node, message, offense),
			}, nil
		} else {
			return nil, nil
		}
	},
}


