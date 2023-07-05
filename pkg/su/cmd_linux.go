package su

import (
	"fmt"
	"os/exec"
	"os/user"
)

const SUName = "root"

func IsSuperUser() bool {
	currentUser, err := user.Current()
	if err != nil {
		return false
	}

	return currentUser.Username == SUName
}

func SUCommand(cmd ...string) (*exec.Cmd, error) {
	if IsSuperUser() {
		return exec.Command(cmd[0], cmd[1:]...), nil
	}

	if !BinaryExists("sudo") {
		return nil, fmt.Errorf("sudo binary not found")
	}

	return exec.Command("sudo", cmd...), nil
}
