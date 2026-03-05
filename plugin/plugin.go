package main

import (
	"github.com/tipo-tawerka/log-linter/internal/analyzer"
	"golang.org/x/tools/go/analysis"
)

type analyzerPlugin struct{}

func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{analyzer.Analyzer}
}

// AnalyzerPlugin — точка входа для golangci-lint.
var AnalyzerPlugin analyzerPlugin

func main() {}
