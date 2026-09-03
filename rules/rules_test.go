package rules_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/TiesDO/gdlint/core"
	"github.com/TiesDO/gdlint/rules"
	"go.lsp.dev/uri"
)

func init() {
	rules.RegisterAll(&core.DefaultRuleRegistry)
}

func loadFixture(name string) ([]byte, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to extract filename from caller stack")
	}

	rootDir := filepath.Dir(filepath.Dir(filename))

	path := filepath.Join(rootDir, "fixtures", "scripts", name)
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func NewDocumentFromFixture(t *testing.T, fixtureName string) *core.Document {
	source, err := loadFixture(fixtureName)

	if err != nil {
		t.Fatal(err)
	}

	document := core.NewDocument(uri.URI(fmt.Sprintf("file:///%s", fixtureName)))
	err = document.UpdateSource(context.Background(), source)

	if err != nil {
		t.Fatal(err)
	}

	return document
}

func NewRunnerWithRule(t *testing.T, rule core.Rule) *core.DocumentRunner {
	runner := core.NewDocumentRunner(&core.DefaultRuleRegistry)
	err := runner.SetRules([]string{rule.Identifier()})

	if err != nil {
		t.Fatal(err)
	}

	return runner
}
