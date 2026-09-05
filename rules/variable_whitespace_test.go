package rules_test

import (
	"testing"

	"github.com/TiesDO/gdlint/core"
	"github.com/TiesDO/gdlint/rules"
	"github.com/stretchr/testify/assert"
)

func TestVariableWhitespaces(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []core.Warning
	}{
		{
			name:     "correct",
			input:    "var b: int = 5",
			expected: []core.Warning{},
		},
		{
			name:  "extra var to name whitespace",
			input: "var   b: int = 5",
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 4,
					EndLine:   0,
					EndChar:   5,
					Message:   "should remove redundant whitespace",
					Offense:   "extra_whitespace",
				},
			},
		},
		{
			name:  "extra name to colon whitespace",
			input: "var b  : int = 5",
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 5,
					EndLine:   0,
					EndChar:   6,
					Message:   "should remove redundant whitespace",
					Offense:   "extra_whitespace",
				},
			},
		},
		{
			name:  "extra colon to type whitespace",
			input: "var b:   int = 5",
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 7,
					EndLine:   0,
					EndChar:   8,
					Message:   "should remove redundant whitespace",
					Offense:   "extra_whitespace",
				},
			},
		},
		{
			name:  "extra type to assign whitespace",
			input: "var b: int   = 5",
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 11,
					EndLine:   0,
					EndChar:   12,
					Message:   "should remove redundant whitespace",
					Offense:   "extra_whitespace",
				},
			},
		},
		{
			name:  "extra assign to value whitespace",
			input: "var b: int =   5",
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 13,
					EndLine:   0,
					EndChar:   14,
					Message:   "should remove redundant whitespace",
					Offense:   "extra_whitespace",
				},
			},
		},
		{
			name:  "extra var to name whitespace inferred type",
			input: "var   b := 5",
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 4,
					EndLine:   0,
					EndChar:   5,
					Message:   "should remove redundant whitespace",
					Offense:   "extra_whitespace",
				},
			},
		},
		{
			name:  "extra name to infer whitespace",
			input: "var b   := 5",
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 6,
					EndLine:   0,
					EndChar:   7,
					Message:   "should remove redundant whitespace",
					Offense:   "extra_whitespace",
				},
			},
		},
		{
			name:  "extra infer to value whitespace",
			input: "var b :=   5",
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 9,
					EndLine:   0,
					EndChar:   10,
					Message:   "should remove redundant whitespace",
					Offense:   "extra_whitespace",
				},
			},
		},
		{
			name:  "extra var to name whitespace no type",
			input: "var   b = 5",
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 4,
					EndLine:   0,
					EndChar:   5,
					Message:   "should remove redundant whitespace",
					Offense:   "extra_whitespace",
				},
			},
		},
		{
			name:  "extra name to assign whitespace no type",
			input: "var b   = 5",
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 6,
					EndLine:   0,
					EndChar:   7,
					Message:   "should remove redundant whitespace",
					Offense:   "extra_whitespace",
				},
			},
		},
		{
			name:  "extra assign to value whitespace no type",
			input: "var b =   5",
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 8,
					EndLine:   0,
					EndChar:   9,
					Message:   "should remove redundant whitespace",
					Offense:   "extra_whitespace",
				},
			},
		},
		{
			name:  "lacking whitespace before assign no type",
			input: "var b= 5",
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 5,
					EndLine:   0,
					EndChar:   6,
					Message:   "should have 1 whitespace character either side",
					Offense:   "no_whitespace",
				},
			},
		},
		{
			name:  "lacking whitespace after assign",
			input: "var b =5",
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 6,
					EndLine:   0,
					EndChar:   7,
					Message:   "should have 1 whitespace character either side",
					Offense:   "no_whitespace",
				},
			},
		},
		{
			name:  "lacking whitespace before and after assign no type",
			input: "var b=5",
			// warning should only appear once
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 5,
					EndLine:   0,
					EndChar:   6,
					Message:   "should have 1 whitespace character either side",
					Offense:   "no_whitespace",
				},
			},
		},
		{
			name:  "lacking whitespace before infer",
			input: "var b:= 5",
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 5,
					EndLine:   0,
					EndChar:   7,
					Message:   "should have 1 whitespace character either side",
					Offense:   "no_whitespace",
				},
			},
		},
		{
			name:  "lacking whitespace after infer",
			input: "var b :=5",
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 6,
					EndLine:   0,
					EndChar:   8,
					Message:   "should have 1 whitespace character either side",
					Offense:   "no_whitespace",
				},
			},
		},
		{
			name:  "lacking whitespace before and after infer",
			input: "var b:=5",
			// warning should only appear once
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 5,
					EndLine:   0,
					EndChar:   7,
					Message:   "should have 1 whitespace character either side",
					Offense:   "no_whitespace",
				},
			},
		},
		{
			name:  "lacking whitespace after colon",
			input: "var b:int = 5",
			expected: []core.Warning{
				{
					StartLine: 0,
					StartChar: 5,
					EndLine:   0,
					EndChar:   6,
					Message:   "should have 1 whitespace character after",
					Offense:   "no_whitespace",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			document := NewDocumentFromString(t, "test.gd", test.input)
			runner := NewRunnerWithRule(t, &rules.VariableWhitespaceRule)

			warnings, err := runner.CheckDocument(document)

			if err != nil {
				t.Fatal(err)
			}

			assert.ElementsMatch(t, test.expected, warnings)
		})
	}
}
