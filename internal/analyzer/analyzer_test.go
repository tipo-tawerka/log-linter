package analyzer

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractStringLit(t *testing.T) {
	tests := []struct {
		name   string
		expr   ast.Expr
		want   string
		wantOK bool
	}{
		{
			"double-quoted string",
			&ast.BasicLit{Kind: token.STRING, Value: "\"hello world\""},
			"hello world",
			true,
		},
		{
			"raw string",
			&ast.BasicLit{Kind: token.STRING, Value: "`raw string`"},
			"raw string",
			true,
		},
		{
			"empty string",
			&ast.BasicLit{Kind: token.STRING, Value: "\"\""},
			"",
			true,
		},
		{
			"int literal",
			&ast.BasicLit{Kind: token.INT, Value: "42"},
			"",
			false,
		},
		{
			"identifier is not a literal",
			&ast.Ident{Name: "someVar"},
			"",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := extractStringLit(tt.expr)
			assert.Equal(t, tt.wantOK, ok)
			if ok {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestParseCSV(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"single value", "lowercase", []string{"lowercase"}},
		{"multiple values", "lowercase,english,specchar", []string{"lowercase", "english", "specchar"}},
		{"with spaces", " lowercase , english ", []string{"lowercase", "english"}},
		{"empty string", "", []string{}},
		{"trailing comma", "lowercase,", []string{"lowercase"}},
		{"only commas", ",,,", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCSV(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAnalyzerDoc(t *testing.T) {
	doc := Analyzer.Doc

	for _, name := range []string{"lowercase", "english", "specchar", "sensitive"} {
		assert.Contains(t, doc, name, "doc should mention rule %q", name)
	}
}

func TestAnalyzerFlags(t *testing.T) {
	f := Analyzer.Flags.Lookup("rules")
	require.NotNil(t, f)
	assert.Contains(t, f.DefValue, "lowercase")

	f = Analyzer.Flags.Lookup("skip-tests")
	require.NotNil(t, f)
	assert.Equal(t, "false", f.DefValue)

	f = Analyzer.Flags.Lookup("sensitive-words")
	require.NotNil(t, f)
	assert.Equal(t, "", f.DefValue)

	f = Analyzer.Flags.Lookup("sensitive-words-file")
	require.NotNil(t, f)
	assert.Equal(t, "", f.DefValue)
}
