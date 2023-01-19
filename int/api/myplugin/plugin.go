package myplugin

import (
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/plugin"
)

func InitializePluginAPI(api *operations.ThyraServerAPI) {
	manager := plugin.NewManager()
	api.PluginManagerInstallHandler = newInstall(manager)
	// api.PluginManagerExecuteCommandHandler
	api.PluginManagerGetInformationHandler = newInfo(manager)
	api.PluginManagerListHandler = newList(manager)
	api.PluginManagerRegisterHandler = newRegister(manager)
	api.PluginManagerUninstallHandler = newUninstall(manager)
}

const (
	errorCodePluginUnknown = "Plugin-0001"

	errorCodePluginInstallationInvalidSource = "Plugin-0010"
)
