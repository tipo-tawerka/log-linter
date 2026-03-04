package rules

import (
	"fmt"
	"regexp"
	"strings"
)

// bannedWords — список запрещенных слов, которые не должны присутствовать в логах.
var bannedWords = []string{
	"password", "secret", "token", "api_key", "apikey",
}

// sensitiveRegexp — регулярное выражение для поиска запрещенных слов в тексте.
var sensitiveRegexp *regexp.Regexp

// AddBannedWords добавляет слова в список запрещенных и пересобирает регулярное выражение.
// Внимание: функция не является потокобезопасной.
func AddBannedWords(words []string) {
	bannedWords = append(bannedWords, words...)
	buildSensitiveRegexp()
}

// buildSensitiveRegexp пересобирает регулярное выражение на основе текущего списка запрещенных слов.
func buildSensitiveRegexp() {
	parts := make([]string, len(bannedWords))
	for i, w := range bannedWords {
		parts[i] = `\b` + regexp.QuoteMeta(w) + `\b`
	}
	sensitiveRegexp = regexp.MustCompile(fmt.Sprintf("(?i)(%s)", strings.Join(parts, "|")))
}

type RuleSensitive struct {
}

func (r RuleSensitive) IsCorrect(text string) bool {
	return !sensitiveRegexp.MatchString(text)
}

func (r RuleSensitive) Fix(text string) string {
	result := sensitiveRegexp.ReplaceAllString(text, "")
	for strings.Contains(result, "  ") {
		result = strings.ReplaceAll(result, "  ", " ")
	}
	return strings.TrimSpace(result)
}

func init() {
	buildSensitiveRegexp()
	addRule(RuleSensitive{}, "sensitive", "требовать отсутствие чувствительных слов")
}
