package lsp_test

import (
	"context"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/TiesDO/gdlint/lsp"
	jrpc "go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	lsp_uri "go.lsp.dev/uri"
)

// testLogWriter bridges a log.Logger to testing.T.Log
type testLogWriter struct {
	t *testing.T
}

func (w testLogWriter) Write(p []byte) (n int, err error) {
	// strings.TrimSpace removes the extra newline log.Printf adds,
	// since t.Log already adds one.
	w.t.Log(strings.TrimSpace(string(p)))
	return len(p), nil
}

type TestClient struct {
	t       *testing.T
	conn    jrpc.Conn
	cancel  context.CancelFunc
	ctx     context.Context
	diagsCh chan protocol.PublishDiagnosticsParams
}

// NewTestClient wires up the pipes, starts the server, and returns the client
func NewTestClient(t *testing.T) *TestClient {
	serverIn, clientOut, _ := os.Pipe()
	clientIn, serverOut, _ := os.Pipe()

	ctx, cancel := context.WithCancel(context.Background())

	logger := log.New(testLogWriter{t: t}, "[SERVER] ", log.Ltime)

	// Start server in the background
	server := lsp.NewServer(serverIn, serverOut, logger)
	go server.Run(ctx)

	diagsCh := make(chan protocol.PublishDiagnosticsParams, 10) // Buffered to prevent blocking

	// Create client
	clientStream := jrpc.NewStream(struct {
		io.Reader
		io.Writer
		io.Closer
	}{
		Reader: clientIn,
		Writer: clientOut,
		Closer: clientIn,
	})

	clientConn := jrpc.NewConn(clientStream)
	clientConn.Go(ctx, func(ctx context.Context, req *jrpc.Request) (any, error) {
		t.Logf("[CLIENT] Received method: %s", req.Method())

		if req.Method() == protocol.MethodTextDocumentPublishDiagnostics {
			var params protocol.PublishDiagnosticsParams

			err := protocol.Unmarshal(req.Params(), &params)
			if err == nil {
				t.Logf("[CLIENT] Successfully unmarshaled diagnostics!")
				diagsCh <- params
			} else {
				t.Errorf("[CLIENT] Failed to unmarshal diagnostics: %v\nRaw JSON: %s", err, string(req.Params()))
			}
		}
		return nil, nil
	})

	return &TestClient{
		t:       t,
		conn:    clientConn,
		cancel:  cancel,
		ctx:     ctx,
		diagsCh: diagsCh,
	}
}

// Close cleans up the context and shuts down the server goroutine
func (tc *TestClient) Close() {
	tc.cancel()
}

// Initialize sends the handshake and returns the server's capabilities
func (tc *TestClient) Initialize() protocol.InitializeResult {
	var result protocol.InitializeResult
	_, err := tc.conn.Call(tc.ctx, protocol.MethodInitialize, nil, &result)
	if err != nil {
		tc.t.Fatalf("Initialize call failed: %v", err)
	}
	return result
}

// DidOpen simulates opening a file in the editor
func (tc *TestClient) DidOpen(uri, text string) {
	params := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        lsp_uri.URI(uri),
			LanguageID: "gdscript",
			Version:    1,
			Text:       text,
		},
	}
	if err := tc.conn.Notify(tc.ctx, protocol.MethodTextDocumentDidOpen, params); err != nil {
		tc.t.Fatalf("DidOpen failed: %v", err)
	}
}

// AwaitDiagnostics blocks until diagnostics arrive or the timeout is reached
func (tc *TestClient) AwaitDiagnostics(timeout time.Duration) protocol.PublishDiagnosticsParams {
	select {
	case diags := <-tc.diagsCh:
		return diags
	case <-time.After(timeout):
		tc.t.Fatalf("Timed out waiting for diagnostics after %v", timeout)
		return protocol.PublishDiagnosticsParams{} // Unreachable due to Fatalf
	}
}
