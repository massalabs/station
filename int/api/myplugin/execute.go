package myplugin

import (
	"fmt"
	"path/filepath"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/plugin"
)

func newExecute(manager *plugin.Manager) operations.PluginManagerExecuteCommandHandler {
	return &execute{manager: manager}
}

type execute struct {
	manager *plugin.Manager
}

//nolint:cyclop
func (c *execute) Handle(params operations.PluginManagerExecuteCommandParams) middleware.Responder {
	plugin, err := c.manager.Plugin(params.ID)
	if err != nil {
		return operations.NewPluginManagerExecuteCommandNotFound().WithPayload(
			&models.Error{Code: errorCodePluginUnknown, Message: fmt.Sprintf("get plugin error: %s", err.Error())})
	}

	status := plugin.Status()

	cmd := params.Body.Command
	switch cmd {
	case "start":
		err := plugin.Start()
		if err != nil {
			return executeFailed(cmd, status,
				fmt.Sprintf("Error while starting plugin %s: %s.\n", filepath.Base(plugin.BinPath), err))
		}
	case "stop":
		err := plugin.Stop()
		if err != nil {
			return executeFailed(cmd, status, fmt.Sprintf("Error while stopping plugin %d: %s.\n", params.ID, err))
		}
	case "restart":
		err := plugin.Stop()
		if err != nil {
			return executeFailed(cmd, status, fmt.Sprintf("Error while stopping plugin %d: %s.\n", params.ID, err))
		}

		err = plugin.Start()
		if err != nil {
			return executeFailed(cmd, status,
				fmt.Sprintf("Error while starting plugin %s: %s.\n", filepath.Base(plugin.BinPath), err))
		}
	case "update":
	default:
		return executeFailed(cmd, status, fmt.Sprintf("Unknown command %s.\n", cmd))
	}

	return operations.NewPluginManagerExecuteCommandNoContent()
}

func executeFailed(cmd string, currentStatus plugin.Status, err string,
) *operations.PluginManagerExecuteCommandBadRequest {
	errStr := ""
	if err != "" {
		errStr = " (err)"
	}

	return operations.NewPluginManagerExecuteCommandBadRequest().WithPayload(
		&models.Error{
			Code:    errorCodePluginExecuteCmdBadRequest,
			Message: fmt.Sprintf("Error: Unable to execute %s command. Current plugin status: %s."+errStr, cmd, currentStatus),
		})
}
