package rules_test

import (
	"context"
	"testing"

	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestUntypedVariableStatements(t *testing.T) {
	runner, err := createDefaultRunnerFromFixture("typed_assignments.gd")

	if err != nil {
		t.Error(err)
	}

	warnings, err := runner.RunRule(rules.UntypedVariableStatementRule.Name(), context.Background())

	if err != nil {
		t.Error(err)
	}

	expected := []rules.Warning{
		{LineNumber: 3, Message: "variable c is untyped", Offense: "untyped_variable_statement"},
		{LineNumber: 8, Message: "variable z is untyped", Offense: "untyped_variable_statement"},
	}

	assert.ElementsMatch(t, expected, warnings)
}
