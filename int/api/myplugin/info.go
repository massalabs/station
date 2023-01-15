package myplugin

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/plugin"
)

//nolint:ireturn
func newInfo(manager *plugin.Manager) operations.PluginManagerGetInformationHandler {
	return &info{manager: manager}
}

type info struct {
	manager *plugin.Manager
}

func (i *info) Handle(param operations.PluginManagerGetInformationParams) middleware.Responder {
	plugin := i.manager.Plugin(param.ID)
	if plugin == nil {
		return operations.NewPluginManagerGetInformationNotFound()
	}

	return operations.NewPluginManagerGetInformationOK().WithPayload(
		&operations.PluginManagerGetInformationOKBody{Status: plugin.Status().String()},
	)
}
