package rules_test

import (
	"testing"

	"github.com/TiesDO/gdlint/core"
	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestCorrectCodeOrder(t *testing.T) {
	document := NewDocumentFromFixture(t, "correct_code_ordering.gd")
	runner := NewRunnerWithRule(t, &rules.CodeOrderRule)

	warnings, err := runner.CheckDocument(document)

	if err != nil {
		t.Error(err)
	}

	assert.Empty(t, warnings)
}

func TestIncorrectCodeOrder(t *testing.T) {
	document := NewDocumentFromFixture(t, "incorrect_code_ordering.gd")
	runner := NewRunnerWithRule(t, &rules.CodeOrderRule)

	warnings, err := runner.CheckDocument(document)

	if err != nil {
		t.Error(err)
	}

	expected := []core.Warning{
		{
			StartLine: 1,
			StartChar: 0,
			EndLine:   1,
			EndChar:   5,
			Message:   "annotations should be defined before class_name_statements",
			Offense:   "code_order",
		},
	}

	assert.ElementsMatch(t, expected, warnings)
}
