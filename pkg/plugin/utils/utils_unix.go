//go:build unix

package utils

func pluginFileName(archiveName string) string {
	return pluginFileNameWithoutExtension(archiveName)
}

func pluginPath(pluginDirectory string, pluginName string) string {
	return pluginPathWithoutExtension(pluginDirectory, pluginName)
}
