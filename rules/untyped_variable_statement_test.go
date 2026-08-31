package rules_test

import (
	"context"
	"testing"

	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestUntypedVariableStatements(t *testing.T) {
	runner, err := createDefaultRunnerFromFixture("typed_assignments.gd")

	if err != nil {
		t.Error(err)
	}

	warnings, err := runner.RunRules([]string{rules.UntypedVariableStatementRule.Name()}, context.Background())

	if err != nil {
		t.Error(err)
	}

	expected := []rules.Warning{
		{
			StartLine: 2,
			StartChar: 4,
			EndLine:   2,
			EndChar:   5,
			Message:   "variable c is untyped",
			Offense:   "untyped_variable_statement",
		}, {
			StartLine: 7,
			StartChar: 5,
			EndLine:   7,
			EndChar:   6,
			Message:   "variable z is untyped",
			Offense:   "untyped_variable_statement",
		},
	}

	assert.ElementsMatch(t, expected, warnings)
}
