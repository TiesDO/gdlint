package core_test

import (
	"context"
	"testing"

	"github.com/TiesDO/gdlint/core"
	"github.com/stretchr/testify/assert"
)

func TestParserString(t *testing.T) {
	parser := core.NewParser()

	err := parser.Parse(context.Background(), []byte("var x = 2"))

	if err != nil {
		t.Fatal(err)
	}

	result := parser.String()
	expected := "source [0, 0] - [0, 9]\n  variable_statement [0, 0] - [0, 9]"

	assert.Contains(t, result, expected)
}

func TestParserQuery(t *testing.T) {
	parser := core.NewParser()

	source := []byte("var x = 2")
	err := parser.Parse(context.Background(), source)

	if err != nil {
		t.Fatal(err)
	}

	matches, err := parser.Query("(variable_statement name: (name) @name)", source)

	if err != nil {
		t.Fatal(err)
	}

	matchCount := 0
	for {
		m := matches.Next()

		if m == nil {
			break
		}

		matchCount++
	}

	assert.Equal(t, matchCount, 1)
}

func TestParserSQuery(t *testing.T) {
	parser := core.NewParser()

	source := []byte("var x = 2")
	err := parser.Parse(context.Background(), source)

	if err != nil {
		t.Fatal(err)
	}

	result, err := parser.SQuery("(variable_statement name: (name) @name)", source)

	if err != nil {
		t.Fatal(err)
	}

	expected := "match: 0\n  capture: name (0)\n  content: x"

	assert.Contains(t, result, expected)
}
