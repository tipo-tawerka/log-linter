package words

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileProvider_Provide(t *testing.T) {
	tests := []struct {
		name string
		data string
		want []string
	}{
		{"with comments and blanks", "password\n# comment\n\nsecret\n  token  \n", []string{"password", "secret", "token"}},
		{"empty file", "", nil},
		{"only comments and blanks", "# comment 1\n\n# comment 2\n", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "words.txt")
			require.NoError(t, os.WriteFile(path, []byte(tt.data), 0644))

			p := FileProvider{}
			got, err := p.Provide(path)
			require.NoError(t, err)
			if tt.want == nil {
				assert.Empty(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestFileProvider_Provide_Error(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{"file not found", "/nonexistent/file.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := FileProvider{}
			_, err := p.Provide(tt.path)
			assert.Error(t, err)
		})
	}
}
