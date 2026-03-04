package lowercase

import (
	"unicode"

	"github.com/tipo-tawerka/log-linter/internal/analyzer/rules"
)

type RuleLowercase struct {
}

func (r RuleLowercase) IsCorrect(text string) bool {
	for _, char := range text {
		if !unicode.IsLetter(char) {
			continue
		}
		if !unicode.IsLower(char) {
			return false
		}
	}
	return true
}

func (r RuleLowercase) Fix(text string) string {
	fixed := make([]rune, 0, len(text))
	for _, char := range text {
		fixed = append(fixed, unicode.ToLower(char))
	}
	return string(fixed)
}

func init() {
	rules.AddRule(RuleLowercase{}, "lowercase", "требовать только строчные буквы")
}
