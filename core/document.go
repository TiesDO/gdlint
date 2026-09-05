package core

import (
	"context"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"go.lsp.dev/uri"
)

type Document struct {
	fileUri uri.URI
	source  []byte
	version int32
	parser  *Parser
	cancel  context.CancelFunc
}

func NewDocument(fileUri uri.URI) *Document {
	d := Document{
		fileUri: fileUri,
		source:  []byte{},
		version: 0,
		parser:  NewParser(),
		cancel:  nil,
	}
	return &d
}

func (d *Document) UpdateSource(ctx context.Context, source []byte) error {
	if d.cancel != nil {
		d.cancel()
	}

	d.source = source
	d.version += 1

	parseCtx, cancel := context.WithCancel(ctx)
	d.cancel = cancel

	return d.parser.Parse(parseCtx, d.source)
}

func (d *Document) Query(pattern string) (*sitter.QueryMatches, *sitter.Query, error) {
	return d.parser.Query(pattern, d.source)
}

func (d *Document) ContentForNode(node *sitter.Node) string {
	return string(d.source[node.StartByte():node.EndByte()])
}

func (d *Document) ByteToPoint(byteOffset uint) sitter.Point {
	if int(byteOffset) >= len(d.source) {
		byteOffset = uint(len(d.source))
	}

	var row uint = 0
	var rowStartByte uint = 0

	for i := uint(0); i < byteOffset; i++ {
		if d.source[i] == '\n' {
			row++
			rowStartByte = i + 1
		}
	}

	return sitter.Point{
		Row:    row,
		Column: byteOffset - rowStartByte,
	}
}
