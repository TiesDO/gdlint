package rules

import (
	_ "embed"
	"fmt"

	"github.com/TiesDO/gdlint/queries"
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

	query, err := sitter.NewQuery(queries.TypeAssignmentPattern, r.language)
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

		for _, capture := range match.Captures {
			node := capture.Node
			startLine := int(node.StartPoint().Row) + 1
			content := string(source[node.StartByte():node.EndByte()])

			captureName := query.CaptureNameForId(capture.Index)
			switch captureName {
			case "untyped_variable_statement":
				warnings = append(warnings, Warning{
					LineNumber: startLine,
					Offense:    "untyped_variable_statement",
					Message:    fmt.Sprintf("variable %s is untyped", content),
				})
			case "untyped_function_return":
				warnings = append(warnings, Warning{
					LineNumber: startLine,
					Offense:    "untyped_function_return",
					Message:    fmt.Sprintf("function %s has no return type", content),
				})
			case "untyped_function_argument":
				warnings = append(warnings, Warning{
					LineNumber: startLine,
					Offense:    "untyped_function_argument",
					Message:    fmt.Sprintf("argument %s is untyped", content),
				})
			}
		}
	}

	return warnings, nil
}
