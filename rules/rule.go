package rules

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
)

type Rule interface {
	Name() string
}

type MatchRule struct {
	name    string
	execute func(*sitter.QueryMatch, []byte) ([]Warning, error)
	pattern []byte
}

func (r *MatchRule) Name() string {
	return r.name
}

func (r *MatchRule) Execute(match *sitter.QueryMatch, source []byte) ([]Warning, error) {
	return r.execute(match, source)
}

type NodeRule struct {
	name    string
	execute func(*sitter.Node, []byte) ([]Warning, error)
}

func (r *NodeRule) Name() string {
	return r.name
}

func (r *NodeRule) Execute(node *sitter.Node, source []byte) ([]Warning, error) {
	return r.execute(node, source)
}
