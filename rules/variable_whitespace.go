package rules

import (
	"errors"

	"github.com/TiesDO/gdlint/core"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

var VariableWhitespaceRule = core.MatchRule{
	Name:    "variable_whitespace",
	Pattern: variableWhitespaceQuery,
	Execute: func(match *sitter.QueryMatch, query *sitter.Query, source []byte) ([]core.Warning, error) {
		if len(match.Captures) < 2 {
			return nil, errors.New("expected at least 2 captures")
		}

		definingName := query.CaptureNames()[match.Captures[1].Index]

		switch definingName {
		case "var.assign":
			name := match.Captures[0]
			assign := match.Captures[1]
			value := match.Captures[2]

			_, e1 := name.Node.ByteRange()
			s2, e2 := assign.Node.ByteRange()
			s3, _ := value.Node.ByteRange()

			nameToAssignDistance := s2 - e1
			assignToValueDistance := s3 - e2

			if nameToAssignDistance == 1 && assignToValueDistance == 1 {
				return nil, nil
			}

			if nameToAssignDistance != 1 && assignToValueDistance != 1 {
				return []core.Warning{
					{
						StartLine: int(name.Node.StartPosition().Row),
						StartChar: int(name.Node.EndPosition().Column),
						EndLine:   int(value.Node.StartPosition().Row),
						EndChar:   int(value.Node.StartPosition().Column),
						Message:   "'=' should have exactly 1 whitespace left and right",
						Offense:   "variable_whitespace",
					},
				}, nil
			} else if nameToAssignDistance != 1 {
				return []core.Warning{
					{
						StartLine: int(name.Node.StartPosition().Row),
						StartChar: int(name.Node.EndPosition().Column),
						EndLine:   int(assign.Node.StartPosition().Row),
						EndChar:   int(assign.Node.StartPosition().Column),
						Message:   "'=' should have exactly 1 whitespace left and right",
						Offense:   "variable_whitespace",
					},
				}, nil
			} else {
				return []core.Warning{
					{
						StartLine: int(assign.Node.StartPosition().Row),
						StartChar: int(assign.Node.EndPosition().Column),
						EndLine:   int(value.Node.StartPosition().Row),
						EndChar:   int(value.Node.StartPosition().Column),
						Message:   "'=' should have exactly 1 whitespace left and right",
						Offense:   "variable_whitespace",
					},
				}, nil
			}
		case "var.colon":
			// handle the explicit typed case
			return nil, errors.New("colon not yet implemented")
		case "var.inferred":
			// handle inferred case
			return nil, errors.New("inferred not yet implemented")
		}

		return nil, nil
	},
}
