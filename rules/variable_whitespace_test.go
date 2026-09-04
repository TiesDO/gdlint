package rules_test

import (
	"testing"

	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestVariableWhitespace(t *testing.T) {
	document := NewDocumentFromFixture(t, "surrounding_whitespace.gd")
	runner := NewRunnerWithRule(t, &rules.VariableWhitespaceRule)

	warnings, err := runner.CheckDocument(document)

	if err != nil {
		t.Fatal(err)
	}

	// expected := []core.Warning{
	// 	{
	// 		StartLine: 4,
	// 		StartChar: 5,
	// 		EndLine:   4,
	// 		EndChar:   6,
	// 		Message:   "'=' should have exactly 1 whitespace left and right",
	// 		Offense:   "variable_whitespace",
	// 	},
	// }

	assert.Len(t, warnings, 4)
	// assert.ElementsMatch(t, expected, warnings)
}
