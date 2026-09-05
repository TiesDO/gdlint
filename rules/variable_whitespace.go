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
			switch captureName {
			case "var.var":
				lastEndByte = capture.Node.EndByte()
			case "var.name":
				startByte := capture.Node.StartByte()
				difference := startByte - lastEndByte

				if difference > 1 {
					warnings = append(warnings, newExtraWhitespaceWarning(lastEndByte+1, startByte-1, document))
				} // else is a never occuring branch

				lastEndByte = capture.Node.EndByte()
			case "var.colon":
				startByte := capture.Node.StartByte()
				difference := startByte - lastEndByte

				if difference > 0 {
					warnings = append(warnings, newExtraWhitespaceWarning(lastEndByte, startByte-1, document))
				}

				lastEndByte = capture.Node.EndByte()
			case "var.type":
				startByte := capture.Node.StartByte()
				difference := startByte - lastEndByte

				if difference > 1 {
					warnings = append(warnings, newExtraWhitespaceWarning(lastEndByte+1, startByte-1, document))
				} else if difference == 0 {
					previous := capture.Node.PrevSibling()
					warnings = append(warnings, newNoWhitespaceWarning(1, previous.StartByte(), previous.EndByte(), document))
				}

				lastEndByte = capture.Node.EndByte()
			case "var.assign":
				startByte := capture.Node.StartByte()
				difference := startByte - lastEndByte

				if difference > 1 {
					warnings = append(warnings, newExtraWhitespaceWarning(lastEndByte+1, startByte-1, document))
				} else if difference == 0 {
					warnings = append(warnings, newNoWhitespaceWarning(0, startByte, capture.Node.EndByte(), document))
				}

				lastEndByte = capture.Node.EndByte()
			case "var.infer":
				startByte := capture.Node.StartByte()
				difference := startByte - lastEndByte

				if difference > 1 {
					warnings = append(warnings, newExtraWhitespaceWarning(lastEndByte+1, startByte-1, document))
				} else if difference == 0 {
					warnings = append(warnings, newNoWhitespaceWarning(0, startByte, capture.Node.EndByte(), document))
				}

				lastEndByte = capture.Node.EndByte()
			case "var.value":
				startByte := capture.Node.StartByte()
				difference := startByte - lastEndByte

				if difference > 1 {
					warnings = append(warnings, newExtraWhitespaceWarning(lastEndByte+1, startByte-1, document))
				} else if difference == 0 {
					previous := capture.Node.PrevSibling()
					warnings = append(warnings, newNoWhitespaceWarning(0, previous.StartByte(), previous.EndByte(), document))
				}

				lastEndByte = capture.Node.EndByte()
			}
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
