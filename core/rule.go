package core

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
)

type Rule interface {
	Identifier() string
}

type MatchRule struct {
	Name    string
	Execute func(*sitter.QueryMatch, []byte) ([]Warning, error)
	Pattern []byte
}

func (r *MatchRule) Identifier() string {
	return r.Name
}

type NodeRule struct {
	Name    string
	Execute func(*sitter.Node, []byte) ([]Warning, error)
}

func (r *NodeRule) Identifier() string {
	return r.Name
}
