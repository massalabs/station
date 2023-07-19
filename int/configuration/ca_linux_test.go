//go:build linux
// +build linux

package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCAPath_Linux(t *testing.T) {
	// Test case for linux with HOME set
	t.Setenv("HOME", "/home/testuser")

	path, err := CAPath()
	assert.NoError(t, err)
	assert.Equal(t, "/home/testuser/.local/share/mkcert", path)
}

func TestCAPath_Linux_XDG_DATA_HOME(t *testing.T) {
	// Test case for linux with XDG_DATA_HOME set
	t.Setenv("XDG_DATA_HOME", "/home/testuser/.local/share")

	path, err := CAPath()
	assert.NoError(t, err)
	assert.Equal(t, "/home/testuser/.local/share/mkcert", path)
}

func TestCAPath_Linux_EmptyHome(t *testing.T) {
	// Test case for linux with empty HOME and XDG_DATA_HOME
	t.Setenv("HOME", "")
	t.Setenv("XDG_DATA_HOME", "")

	_, err := CAPath()
	assert.Error(t, err)
	assert.Equal(t, errMsgCAPath, err.Error())
}
