package myplugin

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/plugin"
)

func info(param operations.PluginManagerGetInformationParams) middleware.Responder {
	if manager == nil {
		manager = plugin.NewManager()
	}

	plugin := manager.Plugin(param.ID)

	return operations.NewPluginManagerGetInformationOK().WithPayload(
		&operations.PluginManagerGetInformationOKBody{Status: plugin.Status().String()},
	)
}
