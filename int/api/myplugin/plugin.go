package myplugin

import (
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/pkg/plugin"
)

func InitializePluginAPI(api *operations.MassastationAPI, pluginManager *plugin.Manager) {
	api.PluginManagerInstallHandler = newInstall(pluginManager)
	api.PluginManagerExecuteCommandHandler = newExecute(pluginManager)
	api.PluginManagerGetInformationHandler = newInfo(pluginManager)
	api.PluginManagerListHandler = newList(pluginManager)
	api.PluginManagerRegisterHandler = newRegister(pluginManager)
	api.PluginManagerUninstallHandler = newUninstall(pluginManager)
	api.PluginManagerGetLogoHandler = newLogo(pluginManager)

	// This endpoint is not defined by the go-swagger API.
	plugin.Handler = *plugin.NewAPIHandler(pluginManager)
}

const (
	errorCodePluginUnknown = "Plugin-0001"

	errorCodePluginInstallationInvalidSource = "Plugin-0010"

	errorCodePluginRegisterUnknown           = "Plugin-0020"
	errorCodePluginRegisterInvalidData       = "Plugin-0020"
	errorCodePluginRegisterAlreadyRegistered = "Plugin-0020"

	errorCodePluginExecuteCmdBadRequest = "Plugin-0030"

	errorCodePluginLogoNotFound = "Plugin-0040"
)
