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
