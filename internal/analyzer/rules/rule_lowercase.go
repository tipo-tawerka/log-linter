package rules

import (
	"unicode"
)

type RuleLowercase struct {
}

func (r RuleLowercase) IsCorrect(text string) bool {
	for _, char := range text {
		if unicode.IsLetter(char) {
			return unicode.IsLower(char)
		}
	}
	return true
}

func (r RuleLowercase) Fix(text string) string {
	runes := []rune(text)
	for i, char := range runes {
		if unicode.IsLetter(char) {
			runes[i] = unicode.ToLower(char)
			break
		}
	}
	return string(runes)
}

func init() {
	addRule(RuleLowercase{}, "lowercase", "требовать только строчные буквы")
}
