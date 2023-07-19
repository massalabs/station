//go:build unix
// +build unix

package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPath_Unix(t *testing.T) {
	path, err := Path()
	assert.NoError(t, err)
	assert.Equal(t, "/usr/local/share/massastation", path)
}
