package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPath_MASSAHOME(t *testing.T) {
	wantedPath := "/test/path"
	// Test case where "MASSA_HOME" is set
	t.Setenv("MASSA_HOME", wantedPath)

	path, err := Path()
	require.NoError(t, err)
	assert.Equal(t, wantedPath, path)
}
