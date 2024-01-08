//go:build unix
// +build unix

package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPath_Unix(t *testing.T) {
	path, err := Path()
	require.NoError(t, err)
	assert.Equal(t, "/usr/local/share/massastation", path)
}
