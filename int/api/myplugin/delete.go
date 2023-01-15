package myplugin

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/plugin"
)

//nolint:ireturn
func newUninstall(manager *plugin.Manager) operations.PluginManagerUninstallHandler {
	return &uninstall{manager: manager}
}

type uninstall struct {
	manager *plugin.Manager
}

func (u *uninstall) Handle(param operations.PluginManagerUninstallParams) middleware.Responder {
	err := u.manager.Delete(param.ID)
	if err != nil {
		return operations.NewPluginManagerUninstallInternalServerError().WithPayload(
			&models.Error{Code: "", Message: err.Error()},
		)
	}

	return operations.NewPluginManagerUninstallNoContent()
}
