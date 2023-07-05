package su

import "os/exec"

func BinaryExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
