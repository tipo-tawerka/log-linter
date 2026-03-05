package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/tipo-tawerka/log-linter/internal/analyzer/rules"
	"github.com/tipo-tawerka/log-linter/internal/analyzer/words"
)

func run(pass *analysis.Pass) (any, error) {
	if err := loadSensitiveWords(); err != nil {
		return nil, err
	}

	activeRules, err := rules.GetRules(parseCSV(flagRules))
	if err != nil {
		return nil, err
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		if shouldSkipTestFile(pass, call.Pos()) {
			return
		}
		msgIndex, ok := logCallMsgIndex(pass, call)
		if !ok || msgIndex >= len(call.Args) {
			return
		}

		msgArg := call.Args[msgIndex]
		if msg, ok := extractStringLit(msgArg); ok {
			checkMessageRules(pass, msgArg, msg, activeRules)
		}
		checkArgsRules(pass, call, activeRules)
	})
	return nil, nil
}

func loadSensitiveWords() error {
	flags := map[string]string{
		words.FlagSensitiveWords: flagSensitiveWords,
		words.FlagSensitiveFile:  flagSensitiveFile,
	}

	for flagName, source := range flags {
		if source == "" {
			continue
		}
		provider, err := words.NewProvider(flagName)
		if err != nil {
			return err
		}
		w, err := provider.Provide(source)
		if err != nil {
			return fmt.Errorf("failed to load sensitive words via %s: %w", flagName, err)
		}
		if len(w) > 0 {
			rules.AddBannedWords(w)
		}
	}
	return nil
}

func shouldSkipTestFile(pass *analysis.Pass, pos token.Pos) bool {
	if !flagSkipTests {
		return false
	}
	file := pass.Fset.File(pos)
	return file != nil && strings.HasSuffix(file.Name(), "_test.go")
}

func checkMessageRules(pass *analysis.Pass, msgArg ast.Expr, msg string, activeRules []rules.RuleData) {
	for _, rd := range activeRules {
		if rd.Rule.IsCorrect(msg) {
			continue
		}

		diag := analysis.Diagnostic{
			Pos:     msgArg.Pos(),
			End:     msgArg.End(),
			Message: fmt.Sprintf("log message violates rule %q: %s", rd.Name, rd.Description),
		}

		if hint, ok := rd.Rule.(rules.CorrectionHint); ok {
			fixed := hint.Fix(msg)
			diag.SuggestedFixes = []analysis.SuggestedFix{{
				Message: fmt.Sprintf("apply rule %q", rd.Name),
				TextEdits: []analysis.TextEdit{{
					Pos:     msgArg.Pos(),
					End:     msgArg.End(),
					NewText: []byte(strconv.Quote(fixed)),
				}}}}
		}
		pass.Report(diag)
	}
}

func checkArgsRules(pass *analysis.Pass, call *ast.CallExpr, activeRules []rules.RuleData) {
	for _, rd := range activeRules {
		checker, ok := rd.Rule.(rules.ArgsChecker)
		if !ok {
			continue
		}
		if checker.CheckArgs(call.Args) {
			continue
		}
		fixHint, hasFixHint := rd.Rule.(rules.ArgsFixHint)
		if !hasFixHint {
			pass.Report(analysis.Diagnostic{
				Pos:     call.Pos(),
				End:     call.End(),
				Message: fmt.Sprintf("log call violates rule %q: %s", rd.Name, rd.Description),
			})
			continue
		}
		for _, v := range fixHint.ArgsViolations(call.Args) {
			pass.Report(analysis.Diagnostic{
				Pos:     v.Pos,
				End:     v.End,
				Message: fmt.Sprintf("variable name %q violates rule %q: %s", v.Ident, rd.Name, rd.Description),
				SuggestedFixes: []analysis.SuggestedFix{{
					Message: fmt.Sprintf("replace %q with %s", v.Ident, rules.RedactedPlaceholder),
					TextEdits: []analysis.TextEdit{{
						Pos:     v.Pos,
						End:     v.End,
						NewText: []byte(rules.RedactedPlaceholder),
					}}}}},
			)
		}
	}
}

func logCallMsgIndex(pass *analysis.Pass, call *ast.CallExpr) (int, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return 0, false
	}

	methodName := sel.Sel.Name

	for _, ld := range registeredLoggers {
		msgIndex, known := ld.Logger.MethodMsgIndex(methodName)
		if !known {
			continue
		}
		if ld.Logger.IsLogCall(pass, sel) {
			return msgIndex, true
		}
	}

	return 0, false
}

func extractStringLit(expr ast.Expr) (string, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}

	val, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}

	return val, true
}

func parseCSV(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
