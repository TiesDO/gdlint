package core

import (
	"context"
	"errors"
	"fmt"
	"strings"

	gdscript "github.com/prestonknopp/tree-sitter-gdscript/bindings/go"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

type Parser struct {
	parser   *sitter.Parser
	language *sitter.Language
	tree     *sitter.Tree
}

func NewParser() *Parser {
	parser := Parser{
		parser: sitter.NewParser(),
	}

	parser.language = sitter.NewLanguage(gdscript.Language())
	parser.parser.SetLanguage(parser.language)

	return &parser
}

func (p *Parser) Parse(ctx context.Context, source []byte) error {
	tree := p.parser.ParseCtx(ctx, source, nil)

	if tree == nil {
		return errors.New("failed to parse tree")
	}

	p.tree = tree

	return nil
}

func (p *Parser) String() string {
	treeCursor := p.tree.Walk()
	var builder strings.Builder

	fprintTreeNode(treeCursor, 0, &builder)

	return builder.String()
}

func (p *Parser) Query(pattern string, source []byte) (*sitter.QueryMatches, error) {
	query, qerr := sitter.NewQuery(p.language, pattern)

	if qerr != nil {
		return nil, fmt.Errorf("failed to create query: %v", qerr.Message)
	}

	cursor := sitter.NewQueryCursor()
	matches := cursor.Matches(query, p.tree.RootNode(), source)

	return &matches, nil
}

func (p *Parser) SQuery(pattern string, source []byte) (string, error) {
	var builder strings.Builder

	query, qerr := sitter.NewQuery(p.language, pattern)

	if qerr != nil {
		return "", fmt.Errorf("failed to create query: %v", qerr.Message)
	}

	cursor := sitter.NewQueryCursor()
	matches := cursor.Matches(query, p.tree.RootNode(), source)

	for {
		match := matches.Next()

		if match == nil {
			break
		}

		fmt.Fprintf(&builder, "match: %d\n", match.Id())
		for _, capture := range match.Captures {
			fmt.Fprintf(&builder, "  capture: %s (%d)\n", capture.Node.GrammarName(), capture.Index)
			fmt.Fprintf(&builder, "  content: %s\n", source[capture.Node.StartByte():capture.Node.EndByte()])
		}
	}

	output := builder.String()

	if len(output) == 0 {
		return "no matches found", nil
	}

	return output, nil
}

func fprintTreeNode(cursor *sitter.TreeCursor, depth int, builder *strings.Builder) {
	indent := strings.Repeat("  ", depth)
	node := cursor.Node()

	prefix := ""
	if fieldName := cursor.FieldName(); fieldName != "" {
		prefix = fieldName + ": "
	}

	start := node.StartPosition()
	end := node.EndPosition()

	fmt.Fprintf(builder,
		"%s%s%s [%d, %d] - [%d, %d]\n",
		indent,
		prefix,
		node.GrammarName(),
		start.Row, start.Column,
		end.Row, end.Column)

	if cursor.GotoFirstChild() {
		fprintTreeNode(cursor, depth+1, builder)
		for cursor.GotoNextSibling() {
			fprintTreeNode(cursor, depth+1, builder)
		}
		cursor.GotoParent()
	}
}
