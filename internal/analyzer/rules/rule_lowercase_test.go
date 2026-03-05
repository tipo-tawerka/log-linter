package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuleLowercase_IsCorrect(t *testing.T) {
	r := RuleLowercase{}

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"lowercase first letter", "starting server", true},
		{"uppercase first letter", "Starting server", false},
		{"all lowercase", "failed to connect", true},
		{"empty string", "", true},
		{"starts with digit", "123 abc", true},
		{"starts with space then uppercase", " Hello", false},
		{"starts with space then lowercase", " hello", true},
		{"only digits", "12345", true},
		{"only punctuation", "!!!", true},
		{"uppercase only", "ERROR", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, r.IsCorrect(tt.input))
		})
	}
}

func TestRuleLowercase_Fix(t *testing.T) {
	r := RuleLowercase{}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"uppercase first letter", "Starting server", "starting server"},
		{"already lowercase", "starting server", "starting server"},
		{"empty string", "", ""},
		{"starts with digit", "123 Abc", "123 abc"},
		{"all caps", "ERROR OCCURRED", "eRROR OCCURRED"},
		{"space then uppercase", " Hello world", " hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, r.Fix(tt.input))
		})
	}
}
