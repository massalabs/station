//go:build windows
// +build windows

package configuration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCAPath_Windows(t *testing.T) {
	// Test case for windows
	path, err := CAPath()
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(os.Getenv("LocalAppData"), "mkcert"), path)
}
