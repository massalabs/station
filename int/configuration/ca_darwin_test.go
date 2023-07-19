//go:build darwin
// +build darwin

package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCAPath_Darwin(t *testing.T) {
	// Test case for darwin with HOME set
	t.Setenv("HOME", "/Users/testuser")

	path, err := CAPath()
	assert.NoError(t, err)
	assert.Equal(t, "/Users/testuser/Library/Application Support/mkcert", path)
}

func TestCAPath_Darwin_EmptyHome(t *testing.T) {
	// Test case for darwin with empty HOME
	t.Setenv("HOME", "")

	_, err := CAPath()
	assert.Error(t, err)
	assert.Equal(t, errMsgCAPath, err.Error())
}
