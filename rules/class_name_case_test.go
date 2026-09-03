package rules_test

import (
	"testing"

	"github.com/TiesDO/gdlint/core"
	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestClassNameCase(t *testing.T) {
	document := NewDocumentFromFixture(t, "class_declarations.gd")
	runner := NewRunnerWithRule(t, &rules.ClassNameCaseRule)

	warnings, err := runner.CheckDocument(document)

	if err != nil {
		t.Error(err)
	}

	expected := []core.Warning{
		{
			StartLine: 0,
			StartChar: 11,
			EndLine:   0,
			EndChar:   14,
			Message:   "class name 'foo' must be PascalCase",
			Offense:   "class_name_case",
		}, {
			StartLine: 2,
			StartChar: 6,
			EndLine:   2,
			EndChar:   9,
			Message:   "class name 'bar' must be PascalCase",
			Offense:   "class_name_case",
		}, {
			StartLine: 8,
			StartChar: 6,
			EndLine:   8,
			EndChar:   14,
			Message:   "class name 'BAR_TEST' must be PascalCase",
			Offense:   "class_name_case",
		},
	}

	assert.ElementsMatch(t, expected, warnings)
}
