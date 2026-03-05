package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuleSpecChar_IsCorrect(t *testing.T) {
	r := RuleSpecChar{}

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"plain words", "server started", true},
		{"letters and digits", "error code 404", true},
		{"exclamation mark", "server started!", false},
		{"multiple exclamation marks", "connection failed!!!", false},
		{"emoji", "server started 🚀", false},
		{"colon", "warning: something went wrong", false},
		{"ellipsis", "loading...", false},
		{"empty string", "", true},
		{"only spaces", "   ", true},
		{"dash", "user-agent not found", false},
		{"slash", "path /api/v1 not found", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, r.IsCorrect(tt.input))
		})
	}
}

func TestRuleSpecChar_Fix(t *testing.T) {
	r := RuleSpecChar{}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"remove exclamation", "server started!", "server started"},
		{"remove multiple specchars", "connection failed!!!", "connection failed"},
		{"remove emoji", "server started 🚀", "server started "},
		{"remove colon", "warning: issue", "warning issue"},
		{"keep letters digits spaces", "error code 404", "error code 404"},
		{"empty string", "", ""},
		{"only special chars", "!@#$%", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, r.Fix(tt.input))
		})
	}
}
