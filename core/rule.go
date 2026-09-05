package core

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
)

type Rule interface {
	Identifier() string
}

type MatchRule struct {
	Name    string
	Execute func(match *sitter.QueryMatch, query *sitter.Query, document *Document) ([]Warning, error)
	Pattern []byte
}

func (r *MatchRule) Identifier() string {
	return r.Name
}

type NodeRule struct {
	Name    string
	Execute func(node *sitter.Node, document *Document) ([]Warning, error)
}

func (r *NodeRule) Identifier() string {
	return r.Name
}
