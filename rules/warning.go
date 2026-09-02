package rules

import (
	"fmt"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

type Warning struct {
	StartLine int
	StartChar int
	EndLine   int
	EndChar   int

	Message string
	Offense string
}

func NewWarningFromNode(node sitter.Node, message string, offense string) *Warning {
	return &Warning{
		StartLine: int(node.StartPosition().Row),
		StartChar: int(node.StartPosition().Column),
		EndLine:   int(node.EndPosition().Row),
		EndChar:   int(node.EndPosition().Column),

		Message: message,
		Offense: offense,
	}
}

func (w *Warning) FullMessage() string {
	return fmt.Sprintf("%d:%d (@%s) - %s", w.StartLine+1, w.StartChar+1, w.Offense, w.Message)
}
