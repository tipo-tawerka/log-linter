package words

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProvider(t *testing.T) {
	tests := []struct {
		name     string
		flag     string
		wantType Provider
	}{
		{"cli provider", FlagSensitiveWords, CLIProvider{}},
		{"file provider", FlagSensitiveFile, FileProvider{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProvider(tt.flag)
			require.NoError(t, err)
			assert.IsType(t, tt.wantType, got)
		})
	}
}

func TestNewProvider_Error(t *testing.T) {
	tests := []struct {
		name string
		flag string
	}{
		{"unknown flag", "bogus"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewProvider(tt.flag)
			assert.Error(t, err)
		})
	}
}
