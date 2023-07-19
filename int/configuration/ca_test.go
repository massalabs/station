package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCAPath_CAROOT(t *testing.T) {
	wantedPath := "/test/path"
	// Test case where "CAROOT" is set
	t.Setenv("CAROOT", wantedPath)

	path, err := CAPath()
	assert.NoError(t, err)
	assert.Equal(t, wantedPath, path)
}
