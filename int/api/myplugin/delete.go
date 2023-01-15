package myplugin

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/plugin"
)

func remove(param operations.PluginManagerUninstallParams) middleware.Responder {
	if manager == nil {
		manager = plugin.NewManager()
	}

	err := manager.Delete(param.ID)
	if err != nil {
		return operations.NewPluginManagerUninstallInternalServerError().WithPayload(
			&models.Error{Code: "", Message: err.Error()},
		)
	}

	return operations.NewPluginManagerUninstallNoContent()
}
