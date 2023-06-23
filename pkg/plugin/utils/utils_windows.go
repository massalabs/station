package utils

import (
	"path/filepath"
	"strings"
)

func pluginFileName(archiveName string) string {
	return strings.Split(archiveName, ".zip")[0] + ".exe"
}

func pluginPath(pluginDirectory string, pluginName string) string {
	return filepath.Join(pluginDirectory, pluginName) + ".exe"
}
