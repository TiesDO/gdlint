package main

import (
	"context"

	"github.com/TiesDO/gdlint/rules"
	gdscript "github.com/prestonknopp/tree-sitter-gdscript/bindings/go"
	sitter "github.com/smacker/go-tree-sitter"
)

var source = []byte("class_name TestClass\nextends Node2D\n\nfunc _ready():\n\tpass\n")

func main() {
	parser := sitter.NewParser()
	language := sitter.NewLanguage(gdscript.Language())
	parser.SetLanguage(language)

	ctx := context.Background()
	tree, err := parser.ParseCtx(ctx, nil, source)

	if err != nil {
		panic("couldn't parse tree")
	}

	rule := rules.NewTypedAssignmentRule(language)
	rule.Check(tree, source)
}
