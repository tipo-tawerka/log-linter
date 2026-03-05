package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// slogMethods — методы slog и индексы аргумента-сообщения.
var slogMethods = map[string]int{
	// Пакетные функции и методы *slog.Logger
	"Debug": 0, "Info": 0, "Warn": 0, "Error": 0,
	// Контекстные варианты
	"DebugContext": 1, "InfoContext": 1, "WarnContext": 1, "ErrorContext": 1,
}

type LoggerSlog struct{}

func (l LoggerSlog) MethodMsgIndex(methodName string) (int, bool) {
	idx, ok := slogMethods[methodName]
	return idx, ok
}

func (l LoggerSlog) IsLogCall(pass *analysis.Pass, sel *ast.SelectorExpr) bool {
	result, ok := l.isPackage(pass, sel)
	if ok {
		return result
	}
	return l.isCall(pass, sel)
}

func (l LoggerSlog) isPackage(pass *analysis.Pass, sel *ast.SelectorExpr) (bool, bool) {
	if ident, ok := sel.X.(*ast.Ident); ok {
		if obj := pass.TypesInfo.Uses[ident]; obj != nil {
			if pkgName, ok := obj.(*types.PkgName); ok {
				return pkgName.Imported().Path() == "log/slog", true
			}
		}
	}
	return false, false
}

func (l LoggerSlog) isCall(pass *analysis.Pass, sel *ast.SelectorExpr) bool {
	typ := pass.TypesInfo.TypeOf(sel.X)
	if typ == nil {
		return false
	}

	typStr := typ.String()
	if ptr, ok := typ.(*types.Pointer); ok {
		typStr = ptr.Elem().String()
	}

	return typStr == "log/slog.Logger"
}

func init() {
	addLogger(LoggerSlog{}, "slog")
}
