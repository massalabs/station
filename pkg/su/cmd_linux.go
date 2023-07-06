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

func NewCommand(cmd ...string) (*Command, error) {
	if IsSuperUser() {
		return &Command{exec.Command(cmd[0], cmd[1:]...)}, nil
	}

	if !BinaryExists("sudo") {
		return nil, fmt.Errorf("sudo binary not found")
	}

	return &Command{exec.Command("sudo", cmd...)}, nil
}
