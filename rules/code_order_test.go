package rules_test

import (
	"context"
	"testing"

	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestCorrectCodeOrder(t *testing.T) {
	runner, err := createDefaultRunnerFromFixture("correct_code_ordering.gd")

	if err != nil {
		t.Fatal(err)
	}

	warnings, err := runner.RunRules([]string{"code_order"}, context.Background())

	if err != nil {
		t.Error(err)
	}

	assert.Empty(t, warnings)
}

func TestIncorrectCodeOrder(t *testing.T) {
	runner, err := createDefaultRunnerFromFixture("incorrect_code_ordering.gd")

	if err != nil {
		t.Fatal(err)
	}

	warnings, err := runner.RunRules([]string{"code_order"}, context.Background())

	if err != nil {
		t.Error(err)
	}

	expected := []rules.Warning{
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
