package sensitive

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tipo-tawerka/log-linter/internal/analyzer/rules"
)

// список запрещенных слов, которые не должны присутствовать в логах.
var bannedWords = []string{
	"password", "secret", "token", "key",
}

// регулярное выражение для поиска запрещенных слов в тексте.
var sensitiveRegexp *regexp.Regexp

// AddBannedWords добавляет слова в список запрещенных слов и сбрасывает регулярное выражение.
// Внимание, функция не является потокобезопасной.
func AddBannedWords(word []string) {
	bannedWords = append(bannedWords, word...)
	sensitiveRegexp = nil
}

type RuleSensitive struct {
}

func (r RuleSensitive) IsCorrect(text string) bool {
	if sensitiveRegexp == nil {
		sensitiveRegexp = regexp.MustCompile(fmt.Sprintf(`(?i)(%s)`, strings.Join(bannedWords, "|")))
	}
	return !sensitiveRegexp.MatchString(text)
}

func (r RuleSensitive) Fix(text string) string {
	if sensitiveRegexp == nil {
		sensitiveRegexp = regexp.MustCompile(fmt.Sprintf(`(?i)(%s)`, strings.Join(bannedWords, "|")))
	}
	return sensitiveRegexp.ReplaceAllString(text, "")
}

func init() {
	rules.AddRule(RuleSensitive{}, "sensitive", "требовать отсутствие чувствительных слов")
}
