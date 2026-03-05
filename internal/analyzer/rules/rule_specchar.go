package rules

import (
	"unicode"
)

type RuleSpecChar struct {
}

func (r RuleSpecChar) IsCorrect(text string) bool {
	for _, char := range text {
		if !r.isCorrectChar(char) {
			return false
		}
	}
	return true
}

func (r RuleSpecChar) Fix(text string) string {
	fixed := make([]rune, 0, len(text))
	for _, char := range text {
		if r.isCorrectChar(char) {
			fixed = append(fixed, char)
		}
	}
	return string(fixed)
}

func (r RuleSpecChar) isCorrectChar(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char) || unicode.IsSpace(char)
}

func init() {
	addRule(RuleSpecChar{}, "specchar", "reuqired without special characters")
}
