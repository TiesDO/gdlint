package rules_test

import (
	"testing"

	"github.com/TiesDO/gdlint/core"
	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestUntypedSignalArgument(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []core.Warning
	}{
		{
			name:   "correct no parameter declaration",
			input:  "signal foo",
			expect: []core.Warning{},
		},
		{
			name:   "correct single parameter declaration",
			input:  "signal foo(a: int)",
			expect: []core.Warning{},
		},
		{
			name:   "correct multi parameter declaration",
			input:  "signal foo(a: int, b: String)",
			expect: []core.Warning{},
		},
		{
			name:  "single parameter declaration",
			input: "signal foo(a)",
			expect: []core.Warning{
				core.Warning{
					StartLine: 0,
					StartChar: 11,
					EndLine:   0,
					EndChar:   12,
					Message:   "untyped argument a",
					Offense:   "untyped_signal_argument",
				},
			},
		},
		{
			name:  "multi parameter declaration",
			input: "signal foo(a, b)",
			expect: []core.Warning{
				core.Warning{
					StartLine: 0,
					StartChar: 11,
					EndLine:   0,
					EndChar:   12,
					Message:   "untyped argument a",
					Offense:   "untyped_signal_argument",
				},
				core.Warning{
					StartLine: 0,
					StartChar: 14,
					EndLine:   0,
					EndChar:   15,
					Message:   "untyped argument b",
					Offense:   "untyped_signal_argument",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			document := NewDocumentFromString(t, "test.gd", test.input)
			runner := NewRunnerWithRule(t, &rules.UntypedSignalArgumentRule)

			warnings, err := runner.CheckDocument(document)

			if err != nil {
				t.Fatal(err)
			}

			assert.ElementsMatch(t, test.expect, warnings)
		})
	}
}
