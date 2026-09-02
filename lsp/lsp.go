package lsp

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

	"github.com/TiesDO/gdlint/rules"
	jrpc "go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
)

const LSP_SOURCE string = "gdlint"

type stdio struct {
	io.Reader
	io.Writer
	io.Closer
}

type Server struct {
	stream jrpc.Stream
	conn   jrpc.Conn
	logger *log.Logger
}

func NewServer(stream_in *os.File, stream_out *os.File, logger *log.Logger) Server {
	if logger == nil {
		logger = log.New(io.Discard, "", 0)
	}

	server := Server{
		logger: logger,
	}

	server.stream = jrpc.NewStream(stdio{
		Reader: stream_in,
		Writer: stream_out,
		Closer: stream_in,
	})

	server.conn = jrpc.NewConn(server.stream)

	return server
}

func (s *Server) Run(ctx context.Context) {
	s.logger.Println("starting listener")
	s.conn.Go(ctx, s.handler)
	<-s.conn.Done()
}

func (s *Server) handler(ctx context.Context, req *jrpc.Request) (any, error) {
	s.logger.Println("handler called with method", req.Method())

	switch req.Method() {
	case protocol.MethodInitialize:
		return s.methodInitialize()
	case protocol.MethodTextDocumentDidOpen:
		return s.methodTextDocumentDidOpen(ctx, req)
	case protocol.MethodTextDocumentDidChange:
		return s.methodTextDocumentDidChange(ctx, req)
	case protocol.MethodShutdown:
		return nil, nil
	}

	return nil, nil
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

	diagnostics, err := s.runDiagnostics(ctx, source)

	if err != nil {
		s.logger.Printf("failed to run diagnostics: %v\n", err)
		return nil, err
	}

	if len(diagnostics) == 0 {
		s.logger.Printf("no diagnostics found\n")
		return diagnostics, nil
	}

	uri := params.TextDocument.URI

	s.logger.Printf("notifying client of %d diagnostics\n", len(diagnostics))

	err = s.conn.Notify(ctx, protocol.MethodTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})

	if err != nil {
		s.logger.Printf("failed to notify diagnostics: %v\n", err)
		return nil, err
	}

	return nil, nil
}

func (s *Server) methodTextDocumentDidChange(ctx context.Context, req *jrpc.Request) (any, error) {
	var params protocol.DidChangeTextDocumentParams

	if err := protocol.Unmarshal(req.Params(), &params); err != nil {
		return nil, err
	}

	if len(params.ContentChanges) == 0 {
		return nil, nil
	}

	full_change_event, ok := params.ContentChanges[0].(*protocol.TextDocumentContentChangeWholeDocument)

	if !ok {
		return nil, errors.New("expected only whole document change events")
	}

	source := []byte(full_change_event.Text)

	diagnostics, err := s.runDiagnostics(ctx, source)

	if err != nil {
		return nil, err
	}

	if len(diagnostics) == 0 {
		s.logger.Printf("no diagnostics found\n")
		return diagnostics, nil
	}

	uri := params.TextDocument.URI

	err = s.conn.Notify(ctx, protocol.MethodTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Server) runDiagnostics(ctx context.Context, source []byte) ([]protocol.Diagnostic, error) {
	runner := rules.NewRuleRunner(&rules.DefaultRuleRegistry, source)
	warnings, err := runner.RunRules(rules.DefaultRuleRegistry.RuleNames(), ctx)

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
