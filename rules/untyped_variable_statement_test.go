package rules_test

import (
	"testing"

	"github.com/TiesDO/gdlint/core"
	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestUntypedVariableStatements(t *testing.T) {
	document := NewDocumentFromFixture(t, "typed_assignments.gd")
	runner := NewRunnerWithRule(t, &rules.UntypedVariableStatementRule)

	warnings, err := runner.CheckDocument(document)

	if err != nil {
		t.Error(err)
	}

	expected := []core.Warning{
		{
			StartLine: 2,
			StartChar: 4,
			EndLine:   2,
			EndChar:   5,
			Message:   "variable c is untyped",
			Offense:   "untyped_variable_statement",
		}, {
			StartLine: 11,
			StartChar: 5,
			EndLine:   11,
			EndChar:   6,
			Message:   "variable z is untyped",
			Offense:   "untyped_variable_statement",
		},
	}

	assert.ElementsMatch(t, expected, warnings)
}
