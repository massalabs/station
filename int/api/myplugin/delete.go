package myplugin

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/plugin"
)

func newUninstall(manager *plugin.Manager) operations.PluginManagerUninstallHandler {
	return &uninstall{manager: manager}
}

type uninstall struct {
	manager *plugin.Manager
}

func (u *uninstall) Handle(param operations.PluginManagerUninstallParams) middleware.Responder {
	logger.Debugf("[DELETE /plugin-manager/%s]", param.ID)

	err := u.manager.Delete(param.ID)
	if err != nil {
		return operations.NewPluginManagerUninstallInternalServerError().WithPayload(
			&models.Error{Code: "", Message: err.Error()},
		)
	}

	return operations.NewPluginManagerUninstallNoContent()
}
