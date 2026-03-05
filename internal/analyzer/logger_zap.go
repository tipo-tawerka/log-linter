package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// zapMethods — методы zap.Logger и zap.SugaredLogger и индексы аргумента-сообщения.
var zapMethods = map[string]int{
	// zap.Logger + zap.SugaredLogger — общие уровни
	"Debug": 0, "Info": 0, "Warn": 0, "Error": 0,
	// zap.Logger — дополнительные уровни
	"DPanic": 0, "Panic": 0, "Fatal": 0,
	// zap.SugaredLogger -f варианты
	"Debugf": 0, "Infof": 0, "Warnf": 0, "Errorf": 0,
	"DPanicf": 0, "Panicf": 0, "Fatalf": 0,
	// zap.SugaredLogger -w варианты
	"Debugw": 0, "Infow": 0, "Warnw": 0, "Errorw": 0,
	"DPanicw": 0, "Panicw": 0, "Fatalw": 0,
}

// zapTypes — типы, распознаваемые как zap-логгеры.
var zapTypes = map[string]bool{
	"go.uber.org/zap.Logger":        true,
	"go.uber.org/zap.SugaredLogger": true,
}

type LoggerZap struct{}

func (l LoggerZap) MethodMsgIndex(methodName string) (int, bool) {
	idx, ok := zapMethods[methodName]
	return idx, ok
}

func (l LoggerZap) IsLogCall(pass *analysis.Pass, sel *ast.SelectorExpr) bool {
	typ := pass.TypesInfo.TypeOf(sel.X)
	if typ == nil {
		return false
	}

	typStr := typ.String()
	if ptr, ok := typ.(*types.Pointer); ok {
		typStr = ptr.Elem().String()
	}

	return zapTypes[typStr]
}

func init() {
	addLogger(LoggerZap{}, "zap")
}
