package lsp_test

import (
	"testing"
	"time"
)

func TestDiagnosticsOnDidOpen(t *testing.T) {
	client := NewTestClient(t)
	defer client.Close()

	// 1. Handshake is usually required by servers before they accept other requests
	client.Initialize()

	// 2. Simulate opening a bad GDScript file
	client.DidOpen("file:///test.gd", "func missing_return():\n  pass")

	// 3. Wait for the server to process and push diagnostics
	diags := client.AwaitDiagnostics(5 * time.Second)

	// 4. Assert
	if diags.URI != "file:///test.gd" {
		t.Errorf("Expected URI file:///test.gd, got %s", diags.URI)
	}

	if len(diags.Diagnostics) != 1 {
		t.Errorf("Expected 1 diagnostic, got %d", len(diags.Diagnostics))
	}
}
