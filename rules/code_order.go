package rules

import (
	"fmt"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

var phaseOrderMap map[string]int = map[string]int{
	"annotation":           1,
	"class_name_statement": 2,
	"extends_statement":    3,
}

var pluralMap map[string]string = map[string]string{
	"annotation":           "annotations",
	"class_name_statement": "class_name_statements",
	"extends_statement":    "extends_statements",
}

var CodeOrderRule = NodeRule{
	name: "code_order",
	execute: func(node *sitter.Node, source []byte) ([]Warning, error) {
		currentPhase := 1
		previousNodeName := "annotation"

		warnings := make([]Warning, 0)

		for i := range node.ChildCount() {
			child := node.Child(uint(i))
			nodeType := child.GrammarName()

			phase, ok := phaseOrderMap[nodeType]

			if !ok {
				// ignore unknown node types for now
				continue
			}

			if phase < currentPhase {
				nodePlural := nodeType
				prevPlural := previousNodeName

				if plural, ok := pluralMap[nodePlural]; ok {
					nodePlural = plural
				}

				if plural, ok := pluralMap[prevPlural]; ok {
					prevPlural = plural
				}

				warnings = append(warnings, Warning{
					StartLine: int(child.StartPosition().Row),
					StartChar: int(child.StartPosition().Column),
					EndLine:   int(child.EndPosition().Row),
					EndChar:   int(child.EndPosition().Column),
					Message:   fmt.Sprintf("%s should be defined before %s", nodePlural, prevPlural),
					Offense:   "code_order",
				})

				currentPhase = phase
				previousNodeName = nodeType
			} else {
				currentPhase = phase
				previousNodeName = nodeType
			}
		}

		return warnings, nil
	},
}

func init() {
	DefaultRuleRegistry.RegisterNodeRule(&CodeOrderRule)
}
