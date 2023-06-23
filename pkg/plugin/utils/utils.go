package utils

func PluginFileName(archiveName string) string {
	return pluginFileName(archiveName)
}

func PluginPath(pluginDirectory string, pluginName string) string {
	return pluginPath(pluginDirectory, pluginName)
}
