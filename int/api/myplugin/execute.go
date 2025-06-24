package myplugin

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/pkg/logger"
	pluginPkg "github.com/massalabs/station/pkg/plugin"
)

func newExecute(manager *pluginPkg.Manager) operations.PluginManagerExecuteCommandHandler {
	return &execute{manager: manager}
}

type execute struct {
	manager *pluginPkg.Manager
}

//nolint:cyclop
func (e *execute) Handle(params operations.PluginManagerExecuteCommandParams) middleware.Responder {
	cmd := params.Body.Command

	logger.Debugf("[POST /plugin-manager/%s/execute] command: %s", params.ID, cmd)

	plugin, err := e.manager.Plugin(params.ID)
	if err != nil {
		return operations.NewPluginManagerExecuteCommandNotFound().WithPayload(
			&models.Error{Code: errorCodePluginUnknown, Message: "get plugin error: " + err.Error()})
	}

	status := plugin.Status()
	pluginName := plugin.Information().Name

	switch cmd {
	case "start":
		err := plugin.Start()
		if err != nil {
			return executeFailed(cmd, status,
				fmt.Sprintf("Error while starting plugin %s: %s.\n", pluginName, err))
		}
	case "stop":
		err := e.manager.StopPlugin(plugin, false)
		if err != nil {
			return executeFailed(cmd, status, fmt.Sprintf("Error while stopping plugin %s: %s.\n", pluginName, err))
		}

	case "restart":
		err := e.manager.StopPlugin(plugin, false)
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

func executeFailed(cmd string, currentStatus pluginPkg.Status, errorMsg string,
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
