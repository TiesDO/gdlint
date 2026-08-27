package rules_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/TiesDO/gdlint/rules"
	gdscript "github.com/prestonknopp/tree-sitter-gdscript/bindings/go"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/stretchr/testify/assert"
)

func fail_on(err error) {
	if err != nil {
		panic(err)
	}
}

func loadFixture(name string) []byte {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to get caller information")
	}
	rootDir := filepath.Dir(filepath.Dir(filename))

	path := filepath.Join(rootDir, "fixtures", "scripts", name)
	data, err := os.ReadFile(path)
	fail_on(err)

	return data
}

func setupTest(fixtureName string) (rules.TypedAssignmentRule, *sitter.Tree, []byte) {
	source := loadFixture(fixtureName)
	parser := sitter.NewParser()
	language := sitter.NewLanguage(gdscript.Language())

	parser.SetLanguage(language)

	ctx := context.Background()
	tree, err := parser.ParseCtx(ctx, nil, source)
	fail_on(err)
	rule := rules.NewTypedAssignmentRule(language)

	return rule, tree, source
}

func TestRootLevelAssignments(t *testing.T) {
	rule, tree, source := setupTest("typed_assignments.gd")
	warnings, err := rule.Check(tree, source)
	fail_on(err)

	expected := []rules.Warning{
		{LineNumber: 3, Message: "variable c is untyped", Offense: "untyped_variable_statement"},
		{LineNumber: 8, Message: "variable z is untyped", Offense: "untyped_variable_statement"},
		{LineNumber: 13, Message: "function bazz has no return type", Offense: "untyped_function_return"},
		{LineNumber: 13, Message: "argument m is untyped", Offense: "untyped_function_argument"},
	}

	assert.ElementsMatch(t, expected, warnings)
}
