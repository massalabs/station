package su

import (
	"fmt"
	"os/exec"
)

type Command struct {
	*exec.Cmd
}

func BinaryExists(name string) bool {
	_, err := exec.LookPath(name)

	return err == nil
}

func (c *Command) Run() error {
	out, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w", string(out), err)
	}

	return nil
}
