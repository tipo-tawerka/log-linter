package rules

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SensitiveSuite struct {
	suite.Suite
	rule          RuleSensitive
	originalWords []string
}

func (s *SensitiveSuite) SetupTest() {
	s.rule = RuleSensitive{}
	s.originalWords = make([]string, len(bannedWords))
	copy(s.originalWords, bannedWords)
}

func (s *SensitiveSuite) TeardownTest() {
	bannedWords = s.originalWords
	buildSensitiveRegexp()
}

func (s *SensitiveSuite) TestIsCorrect_AlwaysTrue() {
	s.True(s.rule.IsCorrect("random text"))
	s.True(s.rule.IsCorrect(""))
}

func (s *SensitiveSuite) TestCheckArgs() {
	tests := []struct {
		name string
		args []ast.Expr
		want bool
	}{
		{
			"clean variable",
			[]ast.Expr{&ast.Ident{Name: "username"}},
			true,
		},
		{
			"exact password",
			[]ast.Expr{&ast.Ident{Name: "password"}},
			false,
		},
		{
			"camelCase userToken",
			[]ast.Expr{&ast.Ident{Name: "userToken"}},
			false,
		},
		{
			"string literal skipped",
			[]ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "\"password\""}},
			true,
		},
		{
			"binary expr with sensitive var",
			[]ast.Expr{&ast.BinaryExpr{
				X:  &ast.BasicLit{Kind: token.STRING, Value: "\"val:\""},
				Op: token.ADD,
				Y:  &ast.Ident{Name: "secretKey"},
			}},
			false,
		},
		{
			"multiple args mixed",
			[]ast.Expr{
				&ast.BasicLit{Kind: token.STRING, Value: "\"user\""},
				&ast.Ident{Name: "username"},
				&ast.BasicLit{Kind: token.STRING, Value: "\"pass\""},
				&ast.Ident{Name: "password"},
			},
			false,
		},
		{
			"all clean",
			[]ast.Expr{
				&ast.Ident{Name: "host"},
				&ast.Ident{Name: "port"},
			},
			true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.Equal(tt.want, s.rule.CheckArgs(tt.args))
		})
	}
}

func (s *SensitiveSuite) TestArgsViolations() {
	tests := []struct {
		name       string
		args       []ast.Expr
		wantIdents []string
	}{
		{
			"no violations",
			[]ast.Expr{&ast.Ident{Name: "username"}},
			nil,
		},
		{
			"single sensitive var",
			[]ast.Expr{&ast.Ident{Name: "password"}},
			[]string{"password"},
		},
		{
			"nested in binary expr",
			[]ast.Expr{&ast.BinaryExpr{
				X:  &ast.BasicLit{Kind: token.STRING, Value: "\"key: \""},
				Op: token.ADD,
				Y:  &ast.Ident{Name: "apikey"},
			}},
			[]string{"apikey"},
		},
		{
			"multiple violations",
			[]ast.Expr{
				&ast.Ident{Name: "password"},
				&ast.Ident{Name: "userToken"},
			},
			[]string{"password", "userToken"},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			violations := s.rule.ArgsViolations(tt.args)
			if len(tt.wantIdents) == 0 {
				s.Empty(violations)
			} else {
				gotIdents := make([]string, len(violations))
				for i, v := range violations {
					gotIdents[i] = v.Ident
				}
				s.Equal(tt.wantIdents, gotIdents)
			}
		})
	}
}

func (s *SensitiveSuite) TestAddBannedWords() {
	s.True(s.rule.CheckArgs([]ast.Expr{&ast.Ident{Name: "creditcard"}}))

	AddBannedWords([]string{"creditcard"})

	s.False(s.rule.CheckArgs([]ast.Expr{&ast.Ident{Name: "creditcard"}}))

	violations := s.rule.ArgsViolations([]ast.Expr{&ast.Ident{Name: "creditcard"}})
	s.Len(violations, 1)
	s.Equal("creditcard", violations[0].Ident)
}

func TestSensitiveSuite(t *testing.T) {
	suite.Run(t, new(SensitiveSuite))
}
