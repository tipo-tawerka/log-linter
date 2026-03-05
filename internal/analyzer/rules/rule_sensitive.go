package rules

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"
)

// bannedWords — список запрещенных слов, которые не должны присутствовать в логах.
var bannedWords = []string{
	"password", "secret", "token", "api_key", "apikey",
}

// sensitiveIdentRegexp — регулярное выражение для поиска запрещенных слов в идентификаторах.
var sensitiveIdentRegexp *regexp.Regexp

// RedactedPlaceholder — замена для чувствительных переменных в исправленном коде.
const RedactedPlaceholder = `"[REDACTED]"`

// AddBannedWords добавляет слова в список запрещенных и пересобирает регулярное выражение.
// Внимание: функция не является потокобезопасной.
func AddBannedWords(words []string) {
	bannedWords = append(bannedWords, words...)
	buildSensitiveRegexp()
}

// buildSensitiveRegexp пересобирает регулярное выражение.
func buildSensitiveRegexp() {
	escaped := make([]string, len(bannedWords))
	for i, w := range bannedWords {
		escaped[i] = regexp.QuoteMeta(w)
	}
	joined := strings.Join(escaped, "|")
	sensitiveIdentRegexp = regexp.MustCompile(fmt.Sprintf(`(?i)(%s)`, joined))
}

// isSensitiveIdent проверяет, содержит ли имя запрещенное слово.
func isSensitiveIdent(name string) bool {
	return sensitiveIdentRegexp.MatchString(name)
}

type RuleSensitive struct {
}

// IsCorrect — заглушка. Проверка чувствительных данных выполняется через ArgsChecker.
func (r RuleSensitive) IsCorrect(_ string) bool {
	return true
}

// CheckArgs проверяет аргументы вызова лог-функции на наличие
// идентификаторов с чувствительными данными. Возвращает true, если все в порядке.
func (r RuleSensitive) CheckArgs(args []ast.Expr) bool {
	for _, arg := range args {
		if hasSensitiveIdent(arg) {
			return false
		}
	}
	return true
}

// ArgsViolations возвращает список нарушений с точными позициями.
func (r RuleSensitive) ArgsViolations(args []ast.Expr) []ArgViolation {
	var violations []ArgViolation
	for _, arg := range args {
		collectViolations(arg, &violations)
	}
	return violations
}

// hasSensitiveIdent рекурсивно проверяет, есть ли в выражении чувствительный идентификатор.
func hasSensitiveIdent(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		return isSensitiveIdent(e.Name)
	case *ast.BinaryExpr:
		return hasSensitiveIdent(e.X) || hasSensitiveIdent(e.Y)
	case *ast.CallExpr:
		for _, a := range e.Args {
			if hasSensitiveIdent(a) {
				return true
			}
		}
	}
	return false
}

// collectViolations рекурсивно извлекает *ast.Ident из выражения
// и собирает нарушения с позициями.
func collectViolations(expr ast.Expr, out *[]ArgViolation) {
	switch e := expr.(type) {
	case *ast.Ident:
		if isSensitiveIdent(e.Name) {
			*out = append(*out, ArgViolation{
				Pos:   e.Pos(),
				End:   e.End(),
				Ident: e.Name,
			})
		}
	case *ast.BinaryExpr:
		collectViolations(e.X, out)
		collectViolations(e.Y, out)
	case *ast.CallExpr:
		for _, a := range e.Args {
			collectViolations(a, out)
		}
	}
}

func init() {
	buildSensitiveRegexp()
	addRule(RuleSensitive{}, "sensitive", "required to avoid logging sensitive data")
}
