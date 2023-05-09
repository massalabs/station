package myplugin

import (
	"fmt"
	"os"

	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/plugin"
)

func InitializePluginAPI(api *operations.ThyraServerAPI) {
	manager, err := plugin.NewManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARN: while starting plugin manager %s.\n", err)
	}

	api.PluginManagerInstallHandler = newInstall(manager)
	api.PluginManagerExecuteCommandHandler = newExecute(manager)
	api.PluginManagerGetInformationHandler = newInfo(manager)
	api.PluginManagerListHandler = newList(manager)
	api.PluginManagerRegisterHandler = newRegister(manager)
	api.PluginManagerUninstallHandler = newUninstall(manager)
	api.PluginManagerUpdateHandler = newUpdate(manager)

	// This endpoint is not defined by the go-swagger API.
	plugin.Handler = *plugin.NewAPIHandler(manager)
}

const (
	errorCodePluginUnknown = "Plugin-0001"

	errorCodePluginInstallationInvalidSource = "Plugin-0010"

	errorCodePluginRegisterUnknown     = "Plugin-0020"
	errorCodePluginRegisterInvalidData = "Plugin-0020"

	errorCodePluginExecuteCmdBadRequest = "Plugin-0030"
)
