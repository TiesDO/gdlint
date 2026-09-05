package rules

import (
	"errors"
	"slices"

	"github.com/TiesDO/gdlint/core"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

var VariableWhitespaceRule = core.MatchRule{
	Name:    "variable_whitespace",
	Pattern: variableWhitespaceQuery,
	Execute: func(match *sitter.QueryMatch, query *sitter.Query, document *core.Document) ([]core.Warning, error) {
		if len(match.Captures) < 1 {
			return nil, errors.New("expected at least 1 capture")
		}

		captureNames := query.CaptureNames()
		var lastEndByte uint = 0
		warnings := []core.Warning{}

		for _, capture := range match.Captures {
			captureName := captureNames[capture.Index]

			if captureName == "var.var" {
				lastEndByte = capture.Node.EndByte()
				continue
			}

			startByte := capture.Node.StartByte()
			difference := startByte - lastEndByte

			switch captureName {
			case "var.name":
				if difference > 1 {
					warnings = append(warnings, newExtraWhitespaceWarning(lastEndByte+1, startByte, document))
				}
				// else is a never occuring branch
			case "var.colon":
				if difference > 0 {
					warnings = append(warnings, newExtraWhitespaceWarning(lastEndByte, startByte, document))
				}
			case "var.assign":
				fallthrough
			case "var.infer":
				if difference > 1 {
					warnings = append(warnings, newExtraWhitespaceWarning(lastEndByte+1, startByte, document))
				} else if difference == 0 {
					warnings = append(warnings, newNoWhitespaceWarning(0, startByte, capture.Node.EndByte(), document))
				}
			case "var.type":
				fallthrough
			case "var.value":
				if difference > 1 {
					warnings = append(warnings, newExtraWhitespaceWarning(lastEndByte+1, startByte, document))
				} else if difference == 0 {
					previous := capture.Node.PrevSibling()

					messageId := 0
					if captureName == "var.type" {
						messageId = 1
					}

					warnings = append(warnings, newNoWhitespaceWarning(messageId, previous.StartByte(), previous.EndByte(), document))
				}
			}

			lastEndByte = capture.Node.EndByte()
		}

		warnings = slices.Compact(warnings)

		return warnings, nil
	},
}

func newExtraWhitespaceWarning(startByte, endByte uint, document *core.Document) core.Warning {
	startPos := document.ByteToPoint(startByte)
	endPos := document.ByteToPoint(endByte)

	return core.Warning{
		StartLine: int(startPos.Row),
		StartChar: int(startPos.Column),
		EndLine:   int(endPos.Row),
		EndChar:   int(endPos.Column),
		Message:   "should remove redundant whitespace",
		Offense:   "extra_whitespace",
	}
}

func newNoWhitespaceWarning(direction int, startByte, endByte uint, document *core.Document) core.Warning {
	startPos := document.ByteToPoint(startByte)
	endPos := document.ByteToPoint(endByte)

	messages := []string{
		"should have 1 whitespace character either side",
		"should have 1 whitespace character after",
	}

	return core.Warning{
		StartLine: int(startPos.Row),
		StartChar: int(startPos.Column),
		EndLine:   int(endPos.Row),
		EndChar:   int(endPos.Column),
		Message:   messages[direction],
		Offense:   "no_whitespace",
	}
}
