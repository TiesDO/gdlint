package lsp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/TiesDO/gdlint/core"
	jrpc "go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

const LSP_SOURCE string = "gdlint"

type stdio struct {
	io.Reader
	io.Writer
	io.Closer
}

type Server struct {
	stream    jrpc.Stream
	conn      jrpc.Conn
	logger    *log.Logger
	documents map[uri.URI]*core.Document
	runner    *core.DocumentRunner
}

func NewServer(stream_in *os.File, stream_out *os.File, logger *log.Logger) Server {
	if logger == nil {
		logger = log.New(io.Discard, "", 0)
	}

	server := Server{
		logger:    logger,
		documents: map[uri.URI]*core.Document{},
		runner:    core.NewDocumentRunner(&core.DefaultRuleRegistry),
	}

	server.stream = jrpc.NewStream(stdio{
		Reader: stream_in,
		Writer: stream_out,
		Closer: stream_in,
	})

	server.conn = jrpc.NewConn(server.stream)

	return server
}

func (s *Server) Run(ctx context.Context) error {
	// TODO: read requested ruleset from a config
	err := s.runner.SetRules(core.DefaultRuleRegistry.RuleNames())

	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	s.logger.Println("starting listener")
	s.conn.Go(ctx, s.handler)
	<-s.conn.Done()
	return nil
}

func (s *Server) handler(ctx context.Context, req *jrpc.Request) (response any, err error) {
	s.logger.Println("handler called with method", req.Method())

	switch req.Method() {
	case protocol.MethodInitialize:
		response, err = s.methodInitialize()
	case protocol.MethodTextDocumentDidOpen:
		response, err = s.methodTextDocumentDidOpen(ctx, req)
	case protocol.MethodTextDocumentDidChange:
		response, err = s.methodTextDocumentDidChange(ctx, req)
	case protocol.MethodShutdown:
		s.logger.Println("executing method called")
		return nil, nil
	}

	if err != nil {
		s.logger.Printf("error while handling method: %v\n", err)
	}

	return response, err
}

func (s *Server) methodInitialize() (any, error) {
	openClose := true
	syncKind := protocol.TextDocumentSyncKindFull

	return protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				OpenClose: &openClose,
				Change:    &syncKind,
			},
		},
	}, nil
}

func (s *Server) methodTextDocumentDidOpen(ctx context.Context, req *jrpc.Request) (any, error) {
	var params protocol.DidOpenTextDocumentParams

	if err := protocol.Unmarshal(req.Params(), &params); err != nil {
		return nil, err
	}

	source := []byte(params.TextDocument.Text)
	uri := params.TextDocument.URI

	document, err := s.ensureDocumentExists(uri, source)
	diagnostics, err := s.runDiagnostics(document)

	if err != nil {
		return nil, err
	}

	s.logger.Printf("notifying client of %d diagnostics\n", len(diagnostics))

	err = s.conn.Notify(ctx, protocol.MethodTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})

	return nil, err
}

func (s *Server) methodTextDocumentDidChange(ctx context.Context, req *jrpc.Request) (any, error) {
	var params protocol.DidChangeTextDocumentParams

	if err := protocol.Unmarshal(req.Params(), &params); err != nil {
		return nil, err
	}

	if len(params.ContentChanges) == 0 {
		s.logger.Printf("no content changes detected, aborting\n")
		return nil, nil
	}

	full_change_event, ok := params.ContentChanges[0].(*protocol.TextDocumentContentChangeWholeDocument)

	if !ok {
		return nil, errors.New("expected only whole document change events")
	}

	source := []byte(full_change_event.Text)
	uri := params.TextDocument.URI

	document, err := s.ensureDocumentExists(uri, source)
	diagnostics, err := s.runDiagnostics(document)

	if err != nil {
		return nil, err
	}

	s.logger.Printf("notifying client of %d diagnostics\n", len(diagnostics))

	err = s.conn.Notify(ctx, protocol.MethodTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})

	return nil, err
}

func (s *Server) runDiagnostics(document *core.Document) ([]protocol.Diagnostic, error) {
	warnings, err := s.runner.CheckDocument(document)

	if err != nil {
		return nil, err
	}

	if len(warnings) == 0 {
		return []protocol.Diagnostic{}, nil
	}

	diagnostics := make([]protocol.Diagnostic, len(warnings))

	for i, warning := range warnings {
		diagnostics[i] = protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(warning.StartLine),
					Character: uint32(warning.StartChar),
				},
				End: protocol.Position{
					Line:      uint32(warning.EndLine),
					Character: uint32(warning.EndChar),
				},
			},
			Code:     protocol.String(warning.Offense),
			Source:   protocol.NewOptional(LSP_SOURCE),
			Message:  protocol.String(warning.Message),
			Severity: protocol.DiagnosticSeverityWarning,
		}
	}

	return diagnostics, nil
}

func (s *Server) ensureDocumentExists(documentUri uri.URI, source []byte) (*core.Document, error) {
	document, ok := s.documents[documentUri]

	if !ok {
		s.logger.Printf("recieved new document %s", documentUri)
		document = core.NewDocument(documentUri)
		s.documents[documentUri] = document
	} else {
		s.logger.Printf("recieved cached document %s", documentUri)
	}

	// figure this out with cancellation tokens so the process doesn't hang
	s.logger.Printf("starting update of source for document %s", documentUri)
	err := document.UpdateSource(context.Background(), source)

	if err != nil {
		return nil, err
	}

	s.logger.Printf("finished update of source for document %s", documentUri)

	return document, nil
}
