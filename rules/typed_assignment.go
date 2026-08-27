package rules

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
)

type TypedAssignmentRule struct {
	language *sitter.Language
}

func NewTypedAssignmentRule(language *sitter.Language) TypedAssignmentRule {
	rule := TypedAssignmentRule{
		language: language,
	}
	return rule
}

func (r *TypedAssignmentRule) Check(tree *sitter.Tree, source []byte) ([]Warning, error) {
	warnings := make([]Warning, 0, 3)
	pattern := []byte("(variable_statement !type name: (name) @var_name) @untyped_assignment")

	query, err := sitter.NewQuery(pattern, r.language)
	if err != nil {
		return warnings, err
	}

	cursor := sitter.NewQueryCursor()

	cursor.Exec(query, tree.RootNode())

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		warning := Warning{
			Offense: "untyped_var",
		}

		for _, capture := range match.Captures {
			captureName := query.CaptureNameForId(capture.Index)

			if captureName == "var_name" {
				node := capture.Node
				varName := string(source[node.StartByte():node.EndByte()])
				warning.Message = fmt.Sprintf("variable %s is untyped", varName)
				warning.LineNumber = int(node.StartPoint().Row) + 1
			}
		}

		warnings = append(warnings, warning)
	}

	return warnings, nil
}
