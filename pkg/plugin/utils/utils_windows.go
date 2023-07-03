package utils

func pluginFileName(archiveName string) string {
	return pluginFileNameWithoutExtension(archiveName) + ".exe"
}

func pluginPath(pluginDirectory string, pluginName string) string {
	return pluginPathWithoutExtension(pluginDirectory, pluginName) + ".exe"
}
