package rules

import (
	sitter "github.com/smacker/go-tree-sitter"
)

type Warning struct {
	LineNumber int
	Message    string
	Offense    string
}

type Rule interface {
	Check(*sitter.Tree, []byte) ([]Warning, error)
}
