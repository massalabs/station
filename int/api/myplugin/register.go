package myplugin

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/plugin"
)

func register(param operations.PluginManagerRegisterParams) middleware.Responder {
	if manager == nil {
		manager = plugin.NewManager()
	}

	wantedPlugin := manager.Plugin(param.Body.ID)
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
