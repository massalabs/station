package myplugin

import (
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
)

func InitializePluginAPI(api *operations.ThyraServerAPI) {
	api.PluginManagerInstallHandler = operations.PluginManagerInstallHandlerFunc(pluginInstall)
	//api.PluginManagerExecuteCommandHandler
	api.PluginManagerGetInformationHandler = operations.PluginManagerGetInformationHandlerFunc(info)
	api.PluginManagerListHandler = operations.PluginManagerListHandlerFunc(list)
	api.PluginManagerRegisterHandler = operations.PluginManagerRegisterHandlerFunc(register)
	api.PluginManagerUninstallHandler = operations.PluginManagerUninstallHandlerFunc(remove)
}

const (
	errorCodePluginUnknown = "Plugin-0001"

	errorCodePluginInstallationInvalidSource = "Plugin-0010"
)
