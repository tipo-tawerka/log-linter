package words

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLIProvider_Provide(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"simple", "password,secret,token", []string{"password", "secret", "token"}},
		{"trim spaces", "  password , secret , token  ", []string{"password", "secret", "token"}},
		{"empty string", "", nil},
		{"trailing comma", "a,b,", []string{"a", "b"}},
		{"single word", "password", []string{"password"}},
		{"only commas", ",,,", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := CLIProvider{}
			got, err := p.Provide(tt.input)
			require.NoError(t, err)
			if tt.want == nil {
				assert.Empty(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
