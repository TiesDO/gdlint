package rules_test

import (
	"testing"

	"github.com/TiesDO/gdlint/core"
	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestUntypedConstStatements(t *testing.T) {
	document := NewDocumentFromFixture(t, "typed_assignments.gd")
	runner := NewRunnerWithRule(t, &rules.UntypedConstStatementRule)

	warnings, err := runner.CheckDocument(document)

	if err != nil {
		t.Error(err)
	}

	expected := []core.Warning{
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
