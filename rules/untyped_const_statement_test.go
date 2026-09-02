package rules_test

import (
	"context"
	"testing"

	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestUntypedConstStatements(t *testing.T) {
	runner, err := createDefaultRunnerFromFixture("typed_assignments.gd")

	if err != nil {
		t.Error(err)
	}

	warnings, err := runner.RunRules([]string{rules.UntypedConstStatementRule.Name()}, context.Background())

	if err != nil {
		t.Error(err)
	}

	expected := []rules.Warning{
		{
			StartLine: 5,
			StartChar: 6,
			EndLine:   5,
			EndChar:   25,
			Message:   "const INFERRED_CONST_TYPE is untyped",
			Offense:   "untyped_const_statement",
		},
		{
			StartLine: 6,
			StartChar: 6,
			EndLine:   6,
			EndChar:   31,
			Message:   "const INVALID_CONST_TYPE_ASSIGN is untyped",
			Offense:   "untyped_const_statement",
		},
	}

	assert.ElementsMatch(t, expected, warnings)
}
