package rules_test

import (
	"context"
	"testing"

	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestUntypedFunctionReturn(t *testing.T) {
	runner, err := createDefaultRunnerFromFixture("typed_assignments.gd")

	if err != nil {
		t.Error(err)
	}

	warnings, err := runner.RunRules([]string{rules.UntypedFunctionReturnRule.Name()}, context.Background())

	if err != nil {
		t.Error(err)
	}

	expected := []rules.Warning{
		{
			StartLine: 12,
			StartChar: 5,
			EndLine:   12,
			EndChar:   9,
			Message:   "function bazz has no return type",
			Offense:   "untyped_function_return",
		},
	}

	assert.ElementsMatch(t, expected, warnings)
}
