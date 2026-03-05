package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuleEnglish_IsCorrect(t *testing.T) {
	r := RuleEnglish{}

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"english only", "starting server on port 8080", true},
		{"russian text", "запуск сервера", false},
		{"mixed english and russian", "server запущен", false},
		{"empty string", "", true},
		{"digits and punctuation", "error 404: not found", true},
		// емодзи не является буквой (unicode.IsLetter возвращает false)
		{"emoji only", "🚀", true},
		{"english with numbers", "failed after 3 retries", true},
		{"chinese characters", "服务器", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, r.IsCorrect(tt.input))
		})
	}
}

func TestRuleEnglish_Fix(t *testing.T) {
	r := RuleEnglish{}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"remove russian letters", "server запущен", "server "},
		{"keep english only", "starting server", "starting server"},
		// Fix удаляет только буквы не из Latin-диапазона. Емодзи — не буква, остаётся.
		{"remove emoji", "server 🚀", "server 🚀"},
		{"remove all non-latin letters", "ошибка подключения", " "},
		{"keep digits and punctuation", "error 404: не найдено", "error 404:  "},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, r.Fix(tt.input))
		})
	}
}
