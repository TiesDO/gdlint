package rules_test

import (
	"testing"

	"github.com/TiesDO/gdlint/core"
	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestUntypedFunctionReturn(t *testing.T) {
	document := NewDocumentFromFixture(t, "typed_assignments.gd")
	runner := NewRunnerWithRule(t, &rules.UntypedFunctionReturnRule)

	warnings, err := runner.CheckDocument(document)

	if err != nil {
		t.Error(err)
	}

	expected := []core.Warning{
		{
			StartLine: 16,
			StartChar: 5,
			EndLine:   16,
			EndChar:   9,
			Message:   "function bazz has no return type",
			Offense:   "untyped_function_return",
		},
	}

	assert.ElementsMatch(t, expected, warnings)
}
