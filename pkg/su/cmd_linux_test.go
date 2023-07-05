package su

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSuperUser(t *testing.T) {
	isSuperUser := IsSuperUser()
	assert.False(t, isSuperUser)
}

func TestSUCommand(t *testing.T) {
	expectedCommand := exec.Command("sudo", "ls", "-l")

	cmd, err := SUCommand("ls", "-l")

	// Check the results
	assert.NoError(t, err)
	assert.Equal(t, cmd.Path, expectedCommand.Path)
	assert.Equal(t, len(cmd.Args), len(expectedCommand.Args))
}
