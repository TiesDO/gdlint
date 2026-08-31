package rules_test

import (
	"context"
	"testing"

	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestConstNameCase(t *testing.T) {
	runner, err := createDefaultRunnerFromFixture("naming_conventions.gd")

	if err != nil {
		t.Error(err)
	}

	warnings, err := runner.RunRules([]string{rules.ConstNameCaseRule.Name()}, context.Background())

	if err != nil {
		t.Error(err)
	}

	expected := []rules.Warning{
		{
			StartLine: 4,
			StartChar: 6,
			EndLine:   4,
			EndChar:   22,
			Message:   "const name 'snake_case_const' must be CONST_CASE",
			Offense:   "const_name_case",
		},
		{
			StartLine: 5,
			StartChar: 6,
			EndLine:   5,
			EndChar:   20,
			Message:   "const name 'camelCaseConst' must be CONST_CASE",
			Offense:   "const_name_case",
		},
		{
			StartLine: 6,
			StartChar: 6,
			EndLine:   6,
			EndChar:   21,
			Message:   "const name 'PascalCaseConst' must be CONST_CASE",
			Offense:   "const_name_case",
		},
		{
			StartLine: 10,
			StartChar: 7,
			EndLine:   10,
			EndChar:   23,
			Message:   "const name 'snake_case_const' must be CONST_CASE",
			Offense:   "const_name_case",
		},
	}

	assert.ElementsMatch(t, expected, warnings)
}
