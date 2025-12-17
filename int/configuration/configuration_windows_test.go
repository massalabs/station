//go:build windows

package configuration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPath_Windows(t *testing.T) {
	expectedPath, err := os.Executable()
	require.NoError(t, err)
	expectedPath = filepath.Dir(expectedPath)

	path, err := Path()
	require.NoError(t, err)
	assert.Equal(t, expectedPath, path)
}
