package myplugin

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/plugin"
)

//nolint:ireturn
func newRegister(manager *plugin.Manager) operations.PluginManagerRegisterHandler {
	return &register{manager: manager}
}

type register struct {
	manager *plugin.Manager
}

func (r *register) Handle(param operations.PluginManagerRegisterParams) middleware.Responder {
	wantedPlugin := r.manager.Plugin(param.Body.ID)
	if wantedPlugin == nil {
		return operations.NewPluginManagerRegisterNotFound()
	}

	info := plugin.Information{
		Name: param.Body.Name, Description: param.Body.Description,
		Logo:      param.Body.Logo,
		Authority: param.Body.Authority, APISpec: param.Body.APISpec,
	}

	wantedPlugin.SetInformation(&info)

	return operations.NewPluginManagerRegisterNoContent()
}
