package rules

import (
	"unicode"
)

type RuleEnglish struct {
}

func (r RuleEnglish) IsCorrect(text string) bool {
	for _, char := range text {
		if !unicode.IsLetter(char) {
			continue
		}
		if !unicode.Is(unicode.Latin, char) {
			return false
		}
	}
	return true
}

func (r RuleEnglish) Fix(text string) string {
	fixed := make([]rune, 0, len(text))
	for _, char := range text {
		if unicode.IsLetter(char) && !unicode.Is(unicode.Latin, char) {
			continue
		}
		fixed = append(fixed, char)
	}
	return string(fixed)
}

func init() {
	addRule(RuleEnglish{}, "english", "требовать только английские буквы")
}
