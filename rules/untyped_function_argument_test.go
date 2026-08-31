package rules_test

import (
	"context"
	"testing"

	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestUntypedFunctionArgument(t *testing.T) {
	runner, err := createDefaultRunnerFromFixture("typed_assignments.gd")

	if err != nil {
		t.Error(err)
	}

	warnings, err := runner.RunRules([]string{rules.UntypedFunctionArgumentRule.Name()}, context.Background())

	if err != nil {
		t.Error(err)
	}

	expected := []rules.Warning{
		{
			StartLine: 12,
			StartChar: 10,
			EndLine:   12,
			EndChar:   11,
			Message:   "untyped argument m",
			Offense:   "untyped_function_argument"},
	}

	assert.ElementsMatch(t, expected, warnings)
}
