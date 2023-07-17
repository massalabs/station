//go:build windows
// +build windows

package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPath_Windows(t *testing.T) {
	expectedPath, err := os.Executable()
	require.NoError(t, err)
	expectedPath = filepath.Dir(expectedPath)

	path, err := Path()
	assert.NoError(t, err)
	assert.Equal(t, expectedPath, path)
}
