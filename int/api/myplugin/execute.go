package myplugin

import (
	"fmt"
	"path/filepath"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/pkg/config"
	"github.com/massalabs/station/pkg/plugin"
)

func newExecute(manager *plugin.Manager) operations.PluginManagerExecuteCommandHandler {
	return &execute{manager: manager}
}

type execute struct {
	manager *plugin.Manager
}

//nolint:cyclop
func (e *execute) Handle(params operations.PluginManagerExecuteCommandParams) middleware.Responder {
	cmd := params.Body.Command

	config.Logger.Debugf("[POST /plugin-manager/%s/execute] command: %s", params.ID, cmd)

	plugin, err := e.manager.Plugin(params.ID)
	if err != nil {
		return operations.NewPluginManagerExecuteCommandNotFound().WithPayload(
			&models.Error{Code: errorCodePluginUnknown, Message: fmt.Sprintf("get plugin error: %s", err.Error())})
	}

	status := plugin.Status()

	pluginName := filepath.Base(plugin.BinPath)

	switch cmd {
	case "start":
		err := plugin.Start()
		if err != nil {
			return executeFailed(cmd, status,
				fmt.Sprintf("Error while starting plugin %s: %s.\n", pluginName, err))
		}
	case "stop":
		err := plugin.Stop()
		if err != nil {
			return executeFailed(cmd, status, fmt.Sprintf("Error while stopping plugin %s: %s.\n", pluginName, err))
		}
	case "restart":
		err := plugin.Stop()
		if err != nil {
			return executeFailed(cmd, status, fmt.Sprintf("Error while stopping plugin %s: %s.\n", pluginName, err))
		}

		err = plugin.Start()
		if err != nil {
			return executeFailed(cmd, status,
				fmt.Sprintf("Error while restarting plugin %s: %s.\n", pluginName, err))
		}
	case "update":
		err := e.manager.Update(params.ID)
		if err != nil {
			return executeFailed(cmd, status, fmt.Sprintf("Error while updating plugin %s: %s.\n", pluginName, err))
		}
	default:
		return executeFailed(cmd, status, fmt.Sprintf("Unknown command %s.\n", cmd))
	}

	return operations.NewPluginManagerExecuteCommandNoContent()
}

func executeFailed(cmd string, currentStatus plugin.Status, errorMsg string,
) *operations.PluginManagerExecuteCommandBadRequest {
	errStr := ""
	if errorMsg != "" {
		errStr = fmt.Sprintf(" (%s)", errorMsg)
	}

	return operations.NewPluginManagerExecuteCommandBadRequest().WithPayload(
		&models.Error{
			Code:    errorCodePluginExecuteCmdBadRequest,
			Message: fmt.Sprintf("[%s] %s. Current plugin status is %s.", cmd, errStr, currentStatus),
		})
}
