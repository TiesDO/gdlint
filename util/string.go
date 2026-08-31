package util

import "unicode"

func IsPascalCase(s string) bool {
	return IsAlphaNumeric(s) && unicode.IsUpper(GetRuneAt(s, 0))
}

func IsConstCase(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}

func IsAlphaNumeric(s string) bool {
	for i := range len(s) {
		if !IsAlphaNumericRune(GetRuneAt(s, i)) {
			return false
		}
	}
	return true
}

func IsAlphaNumericRune(r rune) bool {
	return unicode.IsDigit(r) || unicode.IsUpper(r) || unicode.IsLower(r)
}

func GetRuneAt(str string, i int) rune {
	if len(str) == 0 || i >= len(str) {
		return 0
	}

	rs := []rune(str)
	return rs[i]
}
