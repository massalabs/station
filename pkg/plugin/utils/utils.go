package utils

import (
	"path/filepath"
	"strings"
)

func PluginFileName(archiveName string) string {
	return pluginFileName(archiveName)
}

func pluginFileNameWithoutExtension(archiveName string) string {
	return strings.Split(archiveName, ".zip")[0]
}

func PluginPath(pluginDirectory string, pluginName string) string {
	return pluginPath(pluginDirectory, pluginName)
}

func pluginPathWithoutExtension(pluginDirectory string, pluginName string) string {
	return filepath.Join(pluginDirectory, pluginName)
}
