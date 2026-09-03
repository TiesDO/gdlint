package rules_test

import (
	"testing"

	"github.com/TiesDO/gdlint/core"
	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestConstNameCase(t *testing.T) {
	document := NewDocumentFromFixture(t, "naming_conventions.gd")
	runner := NewRunnerWithRule(t, &rules.ConstNameCaseRule)

	warnings, err := runner.CheckDocument(document)

	if err != nil {
		t.Fatal(err)
	}

	expected := []core.Warning{
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
