package rules_test

import (
	"testing"

	"github.com/TiesDO/gdlint/core"
	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestUntypedFunctionArgument(t *testing.T) {
	document := NewDocumentFromFixture(t, "typed_assignments.gd")
	runner := NewRunnerWithRule(t, &rules.UntypedFunctionArgumentRule)

	warnings, err := runner.CheckDocument(document)

	if err != nil {
		t.Error(err)
	}

	expected := []core.Warning{
		{
			StartLine: 16,
			StartChar: 10,
			EndLine:   16,
			EndChar:   11,
			Message:   "untyped argument m",
			Offense:   "untyped_function_argument"},
	}

	assert.ElementsMatch(t, expected, warnings)
}
