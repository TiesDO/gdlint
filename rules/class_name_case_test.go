package rules_test

import (
	"context"
	"testing"

	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestClassNameCase(t *testing.T) {
	runner, err := createDefaultRunnerFromFixture("class_declarations.gd")

	if err != nil {
		t.Error(err)
	}

	warnings, err := runner.RunRule("class_name_case", context.Background())

	if err != nil {
		t.Error(err)
	}

	expected := []rules.Warning{
		{LineNumber: 1, Message: "class name 'foo' must be PascalCase", Offense: "class_name_case"},
		{LineNumber: 3, Message: "class name 'bar' must be PascalCase", Offense: "class_name_case"},
		{LineNumber: 9, Message: "class name 'BAR_TEST' must be PascalCase", Offense: "class_name_case"},
	}

	assert.ElementsMatch(t, expected, warnings)
}
