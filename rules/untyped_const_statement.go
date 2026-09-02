package rules

import (
	"fmt"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

var UntypedConstStatementRule = MatchRule{
	name:    "untyped_const_statement",
	pattern: []byte("[(const_statement name: (name) @const_name type: (inferred_type)) (const_statement name: (name) @const_name !type)]"),
	execute: func(match *sitter.QueryMatch, source []byte) ([]Warning, error) {
		if len(match.Captures) != 1 {
			return nil, fmt.Errorf("expected only 1 capture, got %d", len(match.Captures))
		}

		capture := match.Captures[0]
		node := capture.Node
		content := string(source[node.StartByte():node.EndByte()])

		warnings := make([]Warning, 1)
		warnings[0] = Warning{
			StartLine: int(node.StartPosition().Row),
			StartChar: int(node.StartPosition().Column),
			EndLine:   int(node.EndPosition().Row),
			EndChar:   int(node.EndPosition().Column),
			Offense:   "untyped_const_statement",
			Message:   fmt.Sprintf("const %s is untyped", content),
		}

		return warnings, nil
	},
}

func init() {
	DefaultRuleRegistry.RegisterMatchRule(&UntypedConstStatementRule)
}
