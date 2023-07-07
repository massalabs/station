package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPath_MASSAHOME(t *testing.T) {
	wantedPath := "/test/path"
	// Test case where "MASSA_HOME" is set
	t.Setenv("MASSA_HOME", wantedPath)

	path, err := Path()
	assert.NoError(t, err)
	assert.Equal(t, wantedPath, path)
}
